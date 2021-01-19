from django.shortcuts import render ,redirect
from  django.http import HttpResponse
from .models import Post
#引入当前的user
from django.contrib.auth.models import User
from django.contrib import messages
from .form import TopicForm,TopicUpdateForm
import requests , json

#当一个class继承了该LoginRequiredMixin 则该class仅在login后才能看
#当一个class继承了该UserPassesTestMixin可设置其中	def test_func(self)方法，为继承该class的类设置使用条件
from django.contrib.auth.mixins import LoginRequiredMixin,UserPassesTestMixin
from django.contrib.auth.decorators import login_required
import http.client
import ssl
import ast
from user.models import Authorization
from user.views import AuthorizationQuerysetIsNotNull
from django.views.generic import ListView ,DetailView,CreateView,UpdateView,DeleteView

url_subsribtion = 'http://localhost:9005/subscriptions'
server_url = "185.128.119.135"

# posts = [
# 	{
# 		'author': 'yaoxin',
# 		'title' : 'blog post 1',
# 		'content': 'first post content',
# 		'data_posted': 'August 27,2018'
# 	},

# 		{
# 		'author': 'yaksdop',
# 		'title' : 'blog post 2',
# 		'content': '2. post content',
# 		'data_posted': 'August 28,2018'
# 	}

# ]


#show all the topic  in homep.html if user is login and subscribed the server
def home(request):
	if request.user.is_authenticated:
		current_user = request.user
		current_profile=current_user.profile

		if current_profile.subscribtionStatus == False:
			messages.warning(request, 'You haven\'t subscribed yet ,please subscribe in profile')
			return redirect('profile')
		else:
			# get all topic from update server

			ulr_gettopic = url_subsribtion+"/"+str(current_profile.subscribtionId)+"/topics"
			#headers={"accept": "application/json"}
			status,data = getAllSubscribtion(request,ulr_gettopic)
			while(status == 401):
				status,data = getAllSubscribtion(request,ulr_gettopic)
			if status == 200:
				messages.success(request, 'Successfully get all topics from update server.')
				#r_list= respons.json() # 将json格式转化为list
				context ={ 'topics':data}
				return render(request,'blog/home.html',context)
			else:
				messages.warning(request, 'Failed to get topic from server')
				return render(request,'blog/home.html')


	else:
		messages.warning(request, 'Please login first')
		return redirect('login')



def about(request):
	return render(request,'blog/about.html')




@login_required	
def topicCreation(request):
	current_user=request.user
	current_profile = current_user.profile 

	if current_profile.subscribtionStatus== False:
		# context ={ 'error_message':"you haven't subscribed yed ,please go to profile to subscribe it"}
		# return render(request,'blog/create_topic.html',context)

		messages.warning(request, 'You haven\'t subscribed yet ,please subscribe in profile')
		return redirect('profile')


	elif request.method == 'POST':
		form = TopicForm(request.POST)

		if form.is_valid():

			threshold = form.cleaned_data['threshold']
			latitude =form.cleaned_data['latitude']
			longitude =form.cleaned_data['longitude']

			date_dic = {
				'threshold':threshold,
				'position' :{ 'latitude':latitude,'longitude':longitude}
				}

			url_topic_creation = "/subscriptions/"+str(current_profile.subscribtionId)+"/topics"
			#headers={"content-type": "application/json","accept": "application/json"} #设置requist 中的传输格式	
			#date= json.dumps(date_dic) # 将dic变为json 格式
			#respons = requests.post(url_topic_creation,date,headers=headers)
			status = sendTopic("POST",request,date_dic,url_topic_creation)
			while(status== 401):
				status = sendTopic("POST",request,date_dic,url_topic_creation)
			if status == 200:
					# topic created
					# delate a id from update server
				# context ={ 'successful_message':"successful creat a topic"}
				messages.success(request, 'successfully creat a topic.')

				# return render(request,'blog/create_topic.html',context)
				return redirect('post-create')

			else:
				# context ={ 'error_message':"creat a topic failed try again"}
				# return render(request,'blog/create_topic.html',context)
				messages.warning(request, 'creat a topic failed.')
				return redirect('post-create')
	else:
		form = TopicForm()

		return render(request,'blog/create_topic.html',{'form': form})


# 				{
#     "threshold": 20,
#     "position": {
#         "latitude": 51.150505,
#         "longitude": 13.763787
#     }
# }

#{{BASE_URL}}/subscriptions/1/topics/1
def detail(request,id):
	current_user=request.user
	current_profile = current_user.profile 
	url_detail = "/subscriptions/"+str(current_profile.subscribtionId)+"/topics/"+str(id)


	#headers={"accept": "application/json"}
	#respons = requests.get(url_detail,headers=headers)
	status,data = getTopic(url_detail)
	while(status==401):
		status,data = getTopic(url_detail)
	
	if status == 200:
		messages.success(request, 'Successfully get the topics from update server.')
		#r_dic= respons.json() # 将json格式转化为dic

		return render(request,'blog/topic_detail.html',{'r_dic':data})
	else:
		messages.warning(request, 'Failed to get this topic ')
		return redirect('blog-home')

#{{BASE_URL}}/subscriptions/1/topics/1
def delate(request,id):
	current_user=request.user
	current_profile = current_user.profile
	url_delate = "/subscriptions/"+str(current_profile.subscribtionId)+"/topics/"+str(id)
	status = deleteTopic(url_delate)
	while(status==401):
		status = deleteTopic(url_delate)
	if status == 204:
		messages.success(request, 'Successful delete the topic from update server.')
		return redirect('blog-home')
	else:
		messages.warning(request, 'Failed to delete this topic, please try again ')
		return redirect('blog-home')
		#return redirect('detail', id=id) # This is the argument of a view


#{{BASE_URL}}/subscriptions/1/topics/1

def update(request,id): 
	current_user=request.user
	current_profile = current_user.profile
	url_update = "/subscriptions/"+str(current_profile.subscribtionId)+"/topics/"+str(id)

	if request.method == 'POST':
		form = TopicUpdateForm(request.POST)

		if form.is_valid():

			threshold = form.cleaned_data['threshold']
			latitude =form.cleaned_data['latitude']
			longitude =form.cleaned_data['longitude']

			data_dic = {
				'threshold':threshold,
				'position' :{ 'latitude':latitude,'longitude':longitude}
				}

			#headers={"content-type": "application/json","accept": "application/json"} #设置requist 中的传输格式
				
			#date= json.dumps(date_dic) # 将dic变为json 格式
			status = sendTopic("PUT",request,data_dic,url_update)
			while(status== 401):
				status = sendTopic("PUT",request,data_dic,url_update)
			if status == 200:
					# topic created
					# delate a id from update server
				# context ={ 'successful_message':"successful creat a topic"}
				messages.success(request, 'Successfully update this topic.')

				# return render(requeslyt,'blog/create_topic.html',context)
				return redirect('blog-home')
			else:

				messages.warning(request, 'Failed to update this topic, please try again')
				return redirect('blog-home')


		else:
			messages.warning(request, 'Incorrect format, please try again')

	form = TopicUpdateForm()

	return render(request,'blog/topic_update.html',{'form': form})


def deleteTopic(url):
	if AuthorizationQuerysetIsNotNull() is False:
		key = getAuthorization()
	else:
		key = Authorization.objects.get(pk=1).authorizationKey
	auth_key = "Bearer "+key
	headers={"content-type": "application/json","accept": "application/hal+json","Authorization":auth_key}
	conn = http.client.HTTPSConnection(server_url,context = ssl._create_unverified_context())
	conn.request("DELETE", url=url, headers=headers)
	response = conn.getresponse()
	if response.status == 204:
		conn.close()
		return 204
	elif response.status == 401:
		getAuthorization()
		conn.close()
		return 401
	else:
		conn.close()
		return 400


def getAllSubscribtion(request,url):
	if AuthorizationQuerysetIsNotNull() is False:
		key = getAuthorization()
	else:
		key = Authorization.objects.get(pk=1).authorizationKey
	auth_key = "Bearer "+key
	headers={"content-type": "application/json","accept": "application/hal+json","Authorization":auth_key}
	conn = http.client.HTTPSConnection(server_url,context = ssl._create_unverified_context())
	conn.request("GET", url=url, headers=headers)
	response = conn.getresponse()
	if response.status==200:
		data = response.read()
		data_str = data.decode("utf-8")
		data_dic = ast.literal_eval(data_str)
		embedded = data_dic['_embedded']
		topics_list=embedded['topics']
		conn.close()
		return 200,topics_list
	elif response.status== 401:
		getAuthorization()
		conn.close()
		return 401,0
	else:
		conn.close()
		return 400,0

def getTopic(url):
	if AuthorizationQuerysetIsNotNull() is False:
		key = getAuthorization()
	else:
		key = Authorization.objects.get(pk=1).authorizationKey	
	auth_key = "Bearer "+key
	headers={"content-type": "application/json","accept": "application/hal+json","Authorization":auth_key}
	conn = http.client.HTTPSConnection(server_url,context = ssl._create_unverified_context())
	conn.request("GET", url=url, headers=headers)
	response = conn.getresponse()
	if response.status == 200:
		data = response.read()
		data_str = data.decode("utf-8")
		data_dic = ast.literal_eval(data_str)
		conn.close()
		return 200,data_dic
	elif response.status== 401:
		getAuthorization()
		conn.close()
		return 401,0
	else:
		conn.close()
		return 400,0

def sendTopic(method,request,date_dic,url):
	if AuthorizationQuerysetIsNotNull() is False:
		key = getAuthorization()
	else:
		key = Authorization.objects.get(pk=1).authorizationKey
	auth_key = "Bearer "+key
	headers={"content-type": "application/json","accept": "application/hal+json","Authorization":auth_key}
	payload_json= json.dumps(date_dic) # 将dic变为json 格式
	conn = http.client.HTTPSConnection(server_url,context = ssl._create_unverified_context())
	conn.request(method=method, url=url, body=payload_json, headers=headers)
	response = conn.getresponse()
	if response.status == 201 or response.status == 200:
		conn.close()
		return 200
	elif response.status== 401:
		getAuthorization()
		conn.close()
		return 401
	else:
		conn.close()
		return 400

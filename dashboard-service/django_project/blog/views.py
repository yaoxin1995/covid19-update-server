from django.shortcuts import render ,redirect
from  django.http import HttpResponse
from .models import Post
#引入当前的user
from django.contrib.auth.models import User
from django.contrib import messages
from .form import TopicForm
import requests , json

#当一个class继承了该LoginRequiredMixin 则该class仅在login后才能看
#当一个class继承了该UserPassesTestMixin可设置其中	def test_func(self)方法，为继承该class的类设置使用条件
from django.contrib.auth.mixins import LoginRequiredMixin,UserPassesTestMixin


from django.views.generic import ListView ,DetailView,CreateView,UpdateView,DeleteView

url_subsribtion = 'http://localhost:9005/subscriptions'

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

	current_user = request.user

	current_profile=current_user.profile

	if current_user.is_authenticated:

		if current_profile.subscribtionStatus == False:
			messages.warning(request, 'you haven\'t subscribed yed ,please subscribe in profile')
			return redirect('profile')
		else:
			# get all topic from update server




	
		context ={ 'posts':current_user.post_set.all()}

		return render(request,'blog/home.html',context)
	else:
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

		messages.warning(request, 'you haven\'t subscribed yed ,please go to profile to subscribe it')
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

			url_topic_creation = url_subsribtion+"/"+str(current_profile.subscribtionId)+"/topics"

			headers={"content-type": "application/json","accept": "application/json"} #设置requist 中的传输格式
				
			date= json.dumps(date_dic) # 将dic变为json 格式
			respons = requests.post(url_topic_creation,date,headers=headers)
			if respons.status_code == 201:
					# topic created
					# delate a id from update server
				# context ={ 'successful_message':"successful creat a topic"}
				messages.success(request, 'successful creat a topic.')

				# return render(request,'blog/create_topic.html',context)
				return redirect('post-create')

			else:
				# context ={ 'error_message':"creat a topic failed try again"}
				# return render(request,'blog/create_topic.html',context)
				messages.success(request, 'successful creat a topic.')
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
from django.shortcuts import render ,redirect

# djago 提供一些 html 的class
from django.contrib.auth.forms import UserCreationForm

from django.contrib import messages

from .form import UserRegistionForm ,UserUpdateForm, ProfileUpdateForm

from django.contrib.auth.decorators import login_required
# Create your views here.

import requests , json



url = 'http://localhost:9005/subscriptions'


def register(request):

#如是post request 则接收网站上的数据
	if request.method == 'POST':
		form = UserRegistionForm(request.POST)

		if form.is_valid():
			#cleaned_data 转换在form中的数据为适合的格式
			#保存form中的内容 密码被自动加密
			form.save()
			username = form.cleaned_data.get('username')
			messages.success(request,f'Account created for {username}，you are now able to login!')
			return redirect('login')
		else:

			messages.warning(request,f'Account creation failed!')
	form = UserRegistionForm()

	#p1:request p2:html的地址 p3:一个要在html中进行查询并显示内容的dictionary
	return render(request,'user/register.html',{'form':form})


@login_required
def profile(request):

	if request.method == 'POST':
		#submit form 
		 # request.user is current login user
		 # request.user.profile is current login user profile
		#  request.POST 是表格中要提交的数据
		u_form =UserUpdateForm(request.POST,instance = request.user)
		p_form = ProfileUpdateForm(request.POST,
								request.FILES, 
								instance= request.user.profile)


		if u_form.is_valid() and p_form.is_valid():

			u_form.save()
			p_form.save()

			subscipted_form = p_form.cleaned_data.get('subscipted')

			current_profile =  request.user.profile

			if subscipted_form != current_profile.subscribtionStatus:
				#get a id from update server
				if subscipted_form == True:

					headers={"content-type": "application/json"} #设置requist 中的传输格式
					date ={ 'email':request.user.email}
					date= json.dumps(date) # 将dic变为json 格式
					respons = requests.post(url,date,headers=headers)
					r_dic= respons.json() # 将json格式转化为dic
					id= r_dic['id']
					current_profile.subscribtionStatus = True
					current_profile.subscribtionId = id 
					current_profile.save()
					# delate a id from update server
				else: 
					current_profile.subscribtionId=0
					current_profile.subscribtionStatus=False
					current_profile.save()  # 204==204



			messages.success(request,f'Account has been updated ')
			return redirect('profile')

	else:
		u_form =UserUpdateForm(instance = request.user)
		p_form = ProfileUpdateForm(instance= request.user.profile)

	content = {
	'u_form': u_form,
	'p_form': p_form
	}
	return render(request,'user/profile.html',content)


# message.debug
# message.info
# message.seccess
# message.warning
# message.error
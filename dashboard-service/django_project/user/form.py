from django import forms
from django.contrib.auth.models import User
from django.contrib.auth.forms import UserCreationForm
from .models import Profile

class UserRegistionForm(UserCreationForm):

	# 向default form 中添加一个 email 选项
	email = forms.EmailField()

	#subscription = forms.BooleanField()

	class Meta:
		#当进行 form.save() 该form 将保存到 User model中
		model = User
		# form 中显示顺序
		fields =['username','email','password1','password2']



class UserUpdateForm(forms.ModelForm):
	class Meta:
		model =User 
		fields =['username','email']


class ProfileUpdateForm(forms.ModelForm):
	class Meta:
		model =Profile 
		#fields =['subscripted','telegram','subscribtionId']
		fields =['subscripted','telegram']


			

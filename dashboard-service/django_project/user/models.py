from django.db import models
from django.contrib.auth.models import User
from PIL import Image
from phonenumber_field.modelfields import PhoneNumberField
# Create your models here.

#不要忘记在当前目录下的admin.py注册model，这样才能在admin下显示该model


class Profile(models.Model):

	# 删除user，profil 也会被删除，反之不会
	user = models.OneToOneField(User,on_delete=models.CASCADE)
	# 每个user均有一个默认的图片，图片被存在 profile_pics文件夹下
	#image = models.ImageField(default ='default.jpg',upload_to = 'profile_pics')

	subscipted = models.BooleanField(default=False)

	phone = PhoneNumberField( blank=False, unique=False,default='000000000000')

	def __str__(self):
		return f'{self.user.username} Profile'

	def save(self, *args, **kwargs):
		super().save(*args, **kwargs)

# #重新定义 该model中的save方法
# 	def save(self,*args, **kwargs):
# 		super().save(*args, **kwargs)

# 		img = Image.open(self.image.path)

# 		if img.height > 300 or img.width>300:
# 			output_sive = (300,300)
# 			img.thumbnail(output_sive)
# 			img.save(self.image.path)


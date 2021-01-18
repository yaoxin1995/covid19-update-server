from django.db import models
from django.utils import timezone
from django.contrib.auth.models import User

from django.urls import reverse

# python manage.py makemigrations
# python manage.py migrate
#得到包含所有user 的QuerySet 
#User.objects.all()
#得到第一个user
#User.objects.first()
#return a QuerySet 里面包含特定的user
#User.objects.filter(username ='xxxx')
#return a 在QuerySet 中的第一个user
#user = User.objects.filter(username ='xxxx').first()

# user.id= user.pk    (primary key)

# Post.objects.all() 得到一个queryset 包含所有post object

#post_1 =Post(titel = 'Blog1 ', content ='First Post Content',author= user)

# post_1.save() //保存post_1 object in databass

# post_1.content 返回 post_1 的content

# post_1.author.email 返回auter 的email

# 得到一个user 的所有 post  格式： .modelname_set 我们可以在这个set上做一些filter的操作
# user.post_set.all() 返回一个query set 包含 该user 所有的post
# user.post_set.creat(title='Blog 3',content ='Third Post content') 创建一个该user的post


# Create your models here.
# 每个class 在这代表数据库中的一个表
class Post(models.Model):
	title = models.CharField(max_length=100)  #character feld
	latitude = models.FloatField(default= 0)#lines lines of text 
	longitude = models.FloatField(default=0)#
	# auto_now = True 每次更新post 即更新时间
	# auto_now_add=True 每次post objec 被创建， 即更新此时间
	# 向其中 传递一个函数 default = timezone.now
	date_posted = models.DateTimeField(default = timezone.now)   

	# USER 是一个表，将user表和该post表用 foreignkey 绑定，on_delete=models.CASCADE
	# 表示当user被删除，则同时删除该post,反之则不。
	author = models.ForeignKey(User,on_delete=models.CASCADE)


	def get_absolute_url(self):
		return reverse('blog-home') 
		# return reverse('post-detail',kwargs={'pk':self.pk})


# class Subscription(models.Model):
# 	subscipted = models.BooleanField(default=False)

# 	author = models.ForeignKey(User,on_delete=models.CASCADE)

	







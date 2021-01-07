from django.contrib import admin
from .models import Profile

# Register your models here.

# python manage.py shell 输入命令行数据

# user= User.objects.filter(username="xxxx").first   得到特定的用户
# user.profile  得到该用户的profile
# user.profile.image 得到其profile的image
# user.profile.image.width/hight....
# user.profile.image.url 图片的location


admin.site.register(Profile)
from django.contrib import admin
from .models import Post
#初始化dadabass 并创建一个superuser
# python manage.py makemigrations
# python manage.py migrate
#python manage.py createsuperuser

# Register your models here.
#在admin page 显示某个model 并对其进行修改
admin.site.register(Post)
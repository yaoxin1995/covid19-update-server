#在 objec 被保存后 该signal 被出发
from django.db.models.signals import post_save
#user model
from django.contrib.auth.models import User 

#reciever is a function that get the signal then perform some task
from django.dispatch import receiver 

from .models import Profile

@receiver(post_save,sender=User)
def create_profile(sender,instance,created,**kwargs):
	if created:
		Profile.objects.create(user=instance)


@receiver(post_save,sender=User)
def save_profile(sender,instance,**kwargs):
		instance.profile.save()

#必须在 apps.py 中注册该信号，才能成功
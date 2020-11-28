from django.shortcuts import render ,redirect
from  django.http import HttpResponse
from .models import Post
#引入当前的user
from django.contrib.auth.models import User


#当一个class继承了该LoginRequiredMixin 则该class仅在login后才能看
#当一个class继承了该UserPassesTestMixin可设置其中	def test_func(self)方法，为继承该class的类设置使用条件
from django.contrib.auth.mixins import LoginRequiredMixin,UserPassesTestMixin


from django.views.generic import ListView ,DetailView,CreateView,UpdateView,DeleteView


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

def home(request):

	current_user = request.user

	if current_user.is_authenticated:
	
		context ={ 'posts':current_user.post_set.all()}

		return render(request,'blog/home.html',context)
	else:
		return redirect('login')



#一个postlist view   home()
class PostListView(ListView):
	"""docstring for ClassName"""
	#查询哪个model 去创建一个list

	# if self.request.user.is_authenticated:
	# 	# model = \
	# 	template_name= 'blog/home.html'   #   <app>/<model>_<viewtype>.html  viewtype 这里指listview
	# 	queryset = current_user.post_set
	# 	context_object_name = 'posts' #指定在 html中list的名字
	# else:
	template_name= 'blog/home.html'   #   <app>/<model>_<viewtype>.html  viewtype 这里指listview
	model = Post
	context_object_name = 'posts'  #指定在 html中list的名字

	#改变 post的出现顺序
	odering = ['data_posted']


#一个postlist view   home()
class PostDetailView(DetailView):
	model= Post



#一个postlist view   home()
class PostCreateView(LoginRequiredMixin,CreateView):
	model= Post
#需要输入的 变量
	fields =['title','latitude','longitude']

	# 给每个提交的form一个auther

	def form_valid(self,form):
		form.instance.author = self.request.user
		return super().form_valid(form)

	# 要想成功创建一个 post 必须在Post model定义其 absolut url def get_absolute_url(self):
	# 这样才能在成功创建后返回


class PostUpdateView(LoginRequiredMixin,UserPassesTestMixin,UpdateView):
	model= Post
#需要输入的 变量
	fields =['title','latitude','longitude']

	# 给每个提交的form一个auther

	def form_valid(self,form):
		form.instance.author = self.request.user
		return super().form_valid(form)

	def test_func(self):
 		post = self.get_object()
 		if self.request.user == post.author:
 			return True
 		else:
 			return False

class PostDeleteView(LoginRequiredMixin,UserPassesTestMixin,DeleteView):
	model= Post

	# delete 后的返回页面
	success_url = '/'

	def test_func(self):
 		post = self.get_object()
 		if self.request.user == post.author:
 			return True
 		else:
 			return False

def about(request):
	return render(request,'blog/about.html',{'title':'about'})
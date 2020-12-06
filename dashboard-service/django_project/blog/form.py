from django import forms




class TopicForm(forms.Form):
	threshold = forms.IntegerField(min_value=0)

	latitude =  forms.IntegerField(min_value=0)

	longitude= forms.IntegerField(min_value=0)
from django import forms




class TopicForm(forms.Form):
	threshold = forms.IntegerField(min_value=0)

	latitude =  forms.FloatField(min_value=0)

	longitude= forms.FloatField(min_value=0)



class TopicUpdateForm(forms.Form):
	threshold = forms.IntegerField(min_value=0)

	latitude =  forms.FloatField(min_value=0)

	longitude= forms.FloatField(min_value=0)
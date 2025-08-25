from django.http import HttpResponse
from django.shortcuts import render
from datetime import datetime

from meetings.models import Meeting


def welcome(request):
    return render(request,"website/welcome.html",
                  {"meetings":Meeting.objects.all(),
                   "x":42})

def date(request):
    return HttpResponse("This page was served at" + str(datetime.now()))

def about(request):
    return HttpResponse("This is a about page")
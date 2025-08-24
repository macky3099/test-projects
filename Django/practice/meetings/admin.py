from django.contrib import admin

from meetings.models import Meeting
from meetings.models import Room

admin.site.register(Meeting)
admin.site.register(Room)
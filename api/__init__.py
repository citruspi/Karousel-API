import ConfigParser

from peewee import *

config = ConfigParser.RawConfigParser()
config.read('server.conf')

database = SqliteDatabase('gallery.db')

from album import AlbumModel
from user import UserModel
from photo import PhotoModel

database.create_tables([PhotoModel, AlbumModel, UserModel], True)
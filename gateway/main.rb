require 'dotenv'
Dotenv.load('.env')

require 'sinatra'
require 'net/http'

require_relative './routes/auth_routes.rb'
require_relative './routes/image_routes.rb'
require_relative './lib/redis.rb'



puts $redis_dd
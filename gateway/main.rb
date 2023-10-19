require 'sinatra'

require_relative './lib/service_discovery.rb'
# require_relative './routes/auth_routes.rb'
require_relative './routes/image_routes.rb'


set :port, 8080

service_discovery = ServiceDiscovery.instance
service_discovery.subscribe
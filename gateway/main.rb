require 'sinatra'

# require_relative './routes/auth_routes.rb'
require_relative './lib/service_discovery.rb'
require_relative './lib/cache_manager.rb'
require_relative './routes/image_routes.rb'
require_relative './routes/config_routes.rb'
require_relative './routes/routes.rb'


set :port, 8080

service_discovery = ServiceDiscovery.instance
service_discovery.subscribe




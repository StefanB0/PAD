# require 'sinatra'

require_relative './routes/auth_routes.rb'
require_relative './routes/image_routes.rb'

require_relative './lib/service_discovery.rb'

set :port, 8080

service_discovery = ServiceDiscovery.instance
service_discovery.subscribe
service_discovery.unsubscribe
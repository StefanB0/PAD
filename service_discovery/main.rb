require 'sinatra'

require_relative 'routes/route.rb'

set :port, 8500

service_db = ServiceDB.instance

Thread.new do
  loop do
    service_db.flush_offline_services
    sleep 1
  end
end
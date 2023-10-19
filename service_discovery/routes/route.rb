require 'sinatra'
require 'json'

require_relative '../lib/service_db.rb'

service_db = ServiceDB.instance

get '/status' do
  [200, 'OK']
end

get '/service/:name' do
  service_arr = service_db.get_service(params[:name])
  if service_arr == :invalid_service
    status 404
  else
    content_type :json
    {services: service_arr}.to_json
  end
end

post '/service' do
  data = JSON.parse(request.body.read)
  name = data["name"]
  address = data["address"]

  if name.nil? || name.nil? || address.empty? || address.empty?
    status 400
    return
  end

  secretkey = service_db.add_service(name, address)
  content_type :json
  [201, {"secret_key" => secretkey}.to_json]
end

delete '/service/:name' do
  return 401 if service_db.remove_service(params[:name], params[:address], params[:secretkey]) == :unauthorized
  status 200
end

get '/status' do
  status [200, 'Discovery Service alive']
end

get '/' do
  404
end

not_found do
  404
end
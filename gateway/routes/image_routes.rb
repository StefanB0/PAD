require 'sinatra'
require 'httparty'
require 'net/http'
require 'json'
require 'logger'


require_relative '../lib/redis.rb'
require_relative '../lib/service_discovery.rb'
require_relative '../lib/load_balancer.rb'

redisCache = RedisCache.instance

# service_discovery = ServiceDiscovery.instance

# img_svc_address = service_discovery.get_service_address('image_service')
# load_balancer = LoadBalancer.new(JSON.parse(img_svc_address)['services'])
load_balancer = LoadBalancer.new('image_service')


logger = Logger.new(STDERR)

get '/image/:id' do
  id = params['id']
  url = URI("#{load_balancer.next_item}/getImage")
  if url.nil?
    return 404
  end

  req = Net::HTTP::Post.new(url, initheader = {'Content-Type' =>'application/json'})
  req.body = {"imageID" => id.to_i}.to_json

  res = Net::HTTP.start(url.hostname, url.port) do |http|
    http.request(req)
  end

  return res.code.to_i if res.code != '200'

  content_type 'image/jpeg'
  res.body
end

get '/image/info/:id' do
  id = params['id']
  url = URI("#{load_balancer.next_item}/getImageInfo")
  

  req = Net::HTTP::Post.new(url, initheader = {'Content-Type' =>'application/json'})
  req.body = {"imageID" => id.to_i}.to_json

  res = Net::HTTP.start(url.hostname, url.port) do |http|
    http.request(req)
  end

  return res.code.to_i if res.code != '200'

  content_type :json
  res.body
end


post '/image' do
  url = "#{load_balancer.next_item}/uploadImage"

  response = HTTParty.post(url, {
    body: {
      token: params[:token],
      author: params[:author],
      title: params[:title],
      description: params[:description],
      tags: params[:tags],
      image: params[:image][:tempfile]
    },
    headers: {
      'Content-Type' => 'multipart/form-data'
    },
    multipart: true
  })

  unless response.code == 201
    logger.error "Error: #{response.code}"
    logger.error response.body.to_s
    return response.code.to_i
  end

  content_type :json
  [201, response.body]
end

post '/image/:id/like' do
  id = params['id']
  url = URI("#{load_balancer.next_item}/likeImage")

  req = Net::HTTP::Post.new(url, initheader = {'Content-Type' =>'application/json'})
  req.body = {"imageID" => id.to_i}.to_json

  res = Net::HTTP.start(url.hostname, url.port) do |http|
    http.request(req)
  end

  return res.code.to_i if res.code != '200'

  content_type :json
  res.body
end


delete '/image/:id' do
  id = params['id']
  url = URI("#{load_balancer.next_item}/deleteImage")

  req = Net::HTTP::Post.new(url, initheader = {'Content-Type' =>'application/json'})
  req.body = {"imageID" => id.to_i}.to_json

  res = Net::HTTP.start(url.hostname, url.port) do |http|
    http.request(req)
  end

  return res.code.to_i if res.code != '200'

  content_type :json
  res.body                
end

put '/image/:id' do
  id = params['id']
  url = URI("#{load_balancer.next_item}/updateImage")

  body = {
    "imageID" => id.to_i,
    "author" => params[:author],
    "title" => params[:title],
    "description" => params[:description]
  }

  req = Net::HTTP::Post.new(url, initheader = {'Content-Type' =>'application/json'})
  req.body = body.to_json

  res = Net::HTTP.start(url.hostname, url.port) do |http|
    http.request(req)
  end

  return res.code.to_i if res.code != '200'

  content_type :json
  res.body                                   
end

# rpc GetImage(GetImageRequest) returns (GetImageResponse) {}
# rpc UploadImage(UploadImageRequest) returns (UploadImageResponse) {}
# rpc DeleteImage(DeleteImageRequest) returns (DeleteImageResponse) {}
# rpc ModifyImage(ModifyImageRequest) returns (ModifyImageResponse) {}
# rpc GetImageList(GetImageListRequest) returns (GetImageListResponse) {}
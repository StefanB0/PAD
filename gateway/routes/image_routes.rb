require 'sinatra'
require 'httparty'
require 'net/http'
require 'json'
require 'logger'


# require_relative '../lib/redis.rb'
require_relative '../lib/service_discovery.rb'
require_relative '../lib/load_balancer.rb'
require_relative '../lib/saga_manager.rb'
require_relative '../lib/circuit_breaker.rb'

# redisCache = RedisCache.instance
sagaManager = SagaManager.instance
cache_manager = CacheManager.instance

# service_discovery = ServiceDiscovery.instance

# img_svc_address = service_discovery.get_service_address('image_service')
# load_balancer = LoadBalancer.new(JSON.parse(img_svc_address)['services'])
load_balancer = LoadBalancer.new('image_service')
circuit_breaker = CircuitBreaker.new('image_service')



logger = Logger.new(STDERR)

get '/image/:id' do
  id = params['id']

  cache = cache_manager.get_key("image:#{id}")

  unless cache.nil?
    logger.info "Cache retrieved for image:#{id}"

    content_type 'image/jpeg'
    return cache
  end

  url = URI("#{load_balancer.next_item}/getImage")
  if url.nil?
    return 404
  end

  req = Net::HTTP::Post.new(url, initheader = {'Content-Type' =>'application/json'})
  req.body = {"imageID" => id.to_i}.to_json

  # res = Net::HTTP.start(url.hostname, url.port) do |http|
  #   http.request(req)
  # end

  res = circuit_breaker.call(req)
  return res if res.is_a? Integer  
  return res.code.to_i if res.code != '200'

  cache_manager.add_key("image:#{id}", res.body)

  content_type 'image/jpeg'
  res.body
end

get '/image/info/:id' do
  id = params['id']
  url = URI("#{load_balancer.next_item}/getImageInfo")
  

  req = Net::HTTP::Post.new(url, initheader = {'Content-Type' =>'application/json'})
  req.body = {"imageID" => id.to_i}.to_json

  res = circuit_breaker.call(req)
  return res if res.is_a? Integer
  return res.code.to_i if res.code != '200'

  content_type :json
  res.body
end


post '/image' do
  # url = "#{load_balancer.next_item}/uploadImage"

  sagaid = sagaManager.new_transaction(['image_service', 'analytics_service'])

  # response = HTTParty.post(url, {
  #   body: {
  #     token: params[:token],
  #     author: params[:author],
  #     title: params[:title],
  #     description: params[:description],
  #     tags: params[:tags],
  #     image: params[:image][:tempfile],
  #     sagaid: sagaid
  #   },
  #   headers: {
  #     'Content-Type' => 'multipart/form-data'
  #   },
  #   multipart: true
  # })

  reqparams = {
    token: params[:token],
    author: params[:author],
    title: params[:title],
    description: params[:description],
    tags: params[:tags],
    image: params[:image][:tempfile],
    sagaid: sagaid
  }
  headers = { 'Content-Type' => 'multipart/form-data' }

  # req = HTTParty::Request.new(Net::HTTP::Post, url, headers: headers, body: reqparams)
  # response = req.perform

  response = circuit_breaker.call_httparty("/uploadImage", Net::HTTP::Post, reqparams, headers)

  return response if response.is_a? Integer

  unless response.code == 201
    logger.error "Error: #{response.code}"
    logger.error response.body.to_s
    return response.code.to_i
  end

  unless sagaManager.check_transaction(sagaid)
    logger.error "Error: Saga failed"
    return 500
  end

  content_type :json
  [201, response.body]
end

post '/image/:id/like' do
  id = params['id']
  url = URI("#{load_balancer.next_item}/likeImage")

  req = Net::HTTP::Post.new(url, initheader = {'Content-Type' =>'application/json'})
  req.body = {"imageID" => id.to_i}.to_json

  res = circuit_breaker.call(req)
  return res if res.is_a? Integer
  return res.code.to_i if res.code != '200'

  content_type :json
  res.body
end


delete '/image/:id' do
  id = params['id']
  url = URI("#{load_balancer.next_item}/deleteImage")

  req = Net::HTTP::Post.new(url, initheader = {'Content-Type' =>'application/json'})
  req.body = {"imageID" => id.to_i}.to_json

  res = circuit_breaker.call(req)
  return res if res.is_a? Integer
  return res.code.to_i if res.code != '200'

  cache_manager.delete_key("image:#{id}")

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

  res = circuit_breaker.call(req)
  return res if res.is_a? Integer
  return res.code.to_i if res.code != '200'

  content_type :json
  res.body                                   
end

put '/transaction/:id' do
  req_body = JSON.parse(request.body.read)
  
  id = params['id']
  service = req_body['service']
  status = req_body['status']

  sagaManager.update_transaction(id, service, status)

  200
end
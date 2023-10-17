require 'sinatra'
require 'json'
require 'grpc'
require 'net/http'
require 'logger'

require_relative '../lib/redis.rb'
require_relative '../lib/service_discovery.rb'

redisCache = RedisCache.instance

service_discovery = ServiceDiscovery.instance

img_svc_address = service_discovery.get_service_address('image_service')

logger = Logger.new(STDERR)

get '/image/:id' do
  id = params['id']
  
  begin
    get_image_request = GetImageRequest.new(id: id)
    get_image_response = nil # todo: call auth service
    
    author = get_image_response.author
    title = get_image_response.title
    description = get_image_response.description
    tags = get_image_response.tags
    imagebytes = get_image_response.imageChunk
    error = get_image_response.error

    if error != 'success'
      return 400
    end

    content_type 'image.png'
    send_data imagebytes, filename: 'image.png', type: 'image/png'
    {"author": author, "title": title, "description": description, "tags": tags}.to_json
  rescue GRPC::BadStatus => e
    logger.error(e.to_json)
    return 400
  end
end

post '/image' do
         
end

delete '/image/:id' do
                
end

put '/image/:id' do
                   
end

get '/allimg' do
                  
end

# rpc GetImage(GetImageRequest) returns (GetImageResponse) {}
# rpc UploadImage(UploadImageRequest) returns (UploadImageResponse) {}
# rpc DeleteImage(DeleteImageRequest) returns (DeleteImageResponse) {}
# rpc ModifyImage(ModifyImageRequest) returns (ModifyImageResponse) {}
# rpc GetImageList(GetImageListRequest) returns (GetImageListResponse) {}
require 'sinatra'
require 'json'
require 'grpc'
require 'net/http'
require 'logger'

require_relative '../lib/redis.rb'
require_relative '../lib/service_discovery.rb'

require_relative '../lib/proto/auth_svc_services_pb.rb'

redisCache = RedisCache.instance

service_discovery = ServiceDiscovery.instance

auth_svc_address = service_discovery.get_service_address('auth_service')
stub = AuthService::Stub.new(auth_svc_address, :this_channel_is_insecure)

logger = Logger.new(STDERR)

post '/login' do
  begin
    login_request = LoginRequest.new(username: params[:username], password: params[:password])
    login_response = stub.Login(login_request)
    
    access_token = login_response.access_token
    refresh_token = login_response.refresh_token

    response.set_cookie('access_token', access_token)
    response.set_cookie('refresh_token', refresh_token)
    200
  rescue GRPC::BadStatus => e
    logger.error(e.to_json)
    return 400
  end
end

post '/signup' do
  begin
    register_request = RegisterRequest.new(username: params[:username], password: params[:password])
    register_response = stub.Register(register_request)
    register_response.error == 'success' ? 200 : 400
  rescue GRPC::BadStatus => e
    logger.error(e.to_json)
    return 400
  end
end

get '/refresh' do
  token = request.cookies['refresh_token']
  unless redisCache.check_token(token)
    return [401, 'Unauthorized']
  end

  begin
    refresh_request = RefreshRequest.new(refresh_token: refresh_token)
    refresh_response = stub.RefreshToken(refresh_request)
    access_token = refresh_response.access_token

    response.set_cookie('access_token', access_token)
    200    
  rescue GRPC::BadStatus => e
    logger.error(e.to_json)
    return 400
  end
end

delete '/user/:id' do
  token = request.cookies['access_token']
  unless redisCache.check_token(token)
    return [401, 'Unauthorized']
  end

  begin
    delete_request = DeleteRequest.new(accesstoken: token)
    delete_response = stub.Delete(delete_request)
    delete_response.error == 'success' ? 200 : 400
  rescue GRPC::BadStatus => e
    logger.error(e.to_json)
    return 400
  end
      
end

get '/allusr' do
  begin
    get_all_response = stub.GetAll(Google::Protobuf::Empty.new)
    get_all_response.users.to_json
  rescue GRPC::BadStatus => e
    logger.error(e.to_json)
    return 400
  end     
end

delete '/allusr' do
  begin
    delete_all_response = stub.DeleteAll(Google::Protobuf::Empty.new)
    delete_all_response.error == 'success' ? 200 : 400
  rescue GRPC::BadStatus => e
    logger.error(e.to_json)
    return 400
  end       
end

# rpc Login(LoginRequest) returns (LoginResponse) {}
# rpc Register(RegisterRequest) returns (RegisterResponse) {}
# rpc RefreshToken(RefreshRequest) returns (RefreshResponse) {}
# rpc Delete(DeleteRequest) returns (DeleteResponse) {}
# rpc GetAll(google.protobuf.Empty) returns (GetAllResponse) {}
# rpc DeleteAll(google.protobuf.Empty) returns (DeleteAllResponse) {}

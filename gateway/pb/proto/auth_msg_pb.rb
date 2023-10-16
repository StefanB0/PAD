# frozen_string_literal: true
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: proto/auth_msg.proto

require 'google/protobuf'


descriptor_data = "\n\x14proto/auth_msg.proto\"2\n\x0cLoginRequest\x12\x10\n\x08username\x18\x01 \x01(\t\x12\x10\n\x08password\x18\x02 \x01(\t\"I\n\rLoginResponse\x12\x13\n\x0b\x61\x63\x63\x65sstoken\x18\x01 \x01(\t\x12\x14\n\x0crefreshtoken\x18\x02 \x01(\t\x12\r\n\x05\x65rror\x18\x03 \x01(\t\"5\n\x0fRegisterRequest\x12\x10\n\x08username\x18\x01 \x01(\t\x12\x10\n\x08password\x18\x02 \x01(\t\"!\n\x10RegisterResponse\x12\r\n\x05\x65rror\x18\x01 \x01(\t\"&\n\x0eRefreshRequest\x12\x14\n\x0crefreshtoken\x18\x01 \x01(\t\"K\n\x0fRefreshResponse\x12\x13\n\x0b\x61\x63\x63\x65sstoken\x18\x01 \x01(\t\x12\x14\n\x0crefreshtoken\x18\x02 \x01(\t\x12\r\n\x05\x65rror\x18\x03 \x01(\t\"$\n\rDeleteRequest\x12\x13\n\x0b\x61\x63\x63\x65sstoken\x18\x01 \x01(\t\"\x1f\n\x0e\x44\x65leteResponse\x12\r\n\x05\x65rror\x18\x01 \x01(\t\".\n\x0eGetAllResponse\x12\r\n\x05users\x18\x01 \x03(\t\x12\r\n\x05\x65rror\x18\x02 \x01(\t\"\"\n\x11\x44\x65leteAllResponse\x12\r\n\x05\x65rror\x18\x01 \x01(\tb\x06proto3"

pool = Google::Protobuf::DescriptorPool.generated_pool

begin
  pool.add_serialized_file(descriptor_data)
rescue TypeError => e
  # Compatibility code: will be removed in the next major version.
  require 'google/protobuf/descriptor_pb'
  parsed = Google::Protobuf::FileDescriptorProto.decode(descriptor_data)
  parsed.clear_dependency
  serialized = parsed.class.encode(parsed)
  file = pool.add_serialized_file(serialized)
  warn "Warning: Protobuf detected an import path issue while loading generated file #{__FILE__}"
  imports = [
  ]
  imports.each do |type_name, expected_filename|
    import_file = pool.lookup(type_name).file_descriptor
    if import_file.name != expected_filename
      warn "- #{file.name} imports #{expected_filename}, but that import was loaded as #{import_file.name}"
    end
  end
  warn "Each proto file must use a consistent fully-qualified name."
  warn "This will become an error in the next major version."
end

LoginRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("LoginRequest").msgclass
LoginResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("LoginResponse").msgclass
RegisterRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("RegisterRequest").msgclass
RegisterResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("RegisterResponse").msgclass
RefreshRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("RefreshRequest").msgclass
RefreshResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("RefreshResponse").msgclass
DeleteRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("DeleteRequest").msgclass
DeleteResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("DeleteResponse").msgclass
GetAllResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("GetAllResponse").msgclass
DeleteAllResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("DeleteAllResponse").msgclass

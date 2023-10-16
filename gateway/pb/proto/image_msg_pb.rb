# frozen_string_literal: true
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: proto/image_msg.proto

require 'google/protobuf'


descriptor_data = "\n\x15proto/image_msg.proto\"y\n\x12UploadImageRequest\x12\r\n\x05token\x18\x01 \x01(\t\x12\x0e\n\x06\x61uthor\x18\x02 \x01(\t\x12\r\n\x05title\x18\x03 \x01(\t\x12\x13\n\x0b\x64\x65scription\x18\x04 \x01(\t\x12\x0c\n\x04Tags\x18\x05 \x03(\t\x12\x12\n\nimageChunk\x18\x06 \x01(\x0c\"5\n\x13UploadImageResponse\x12\x0f\n\x07imageID\x18\x01 \x01(\x03\x12\r\n\x05\x65rror\x18\x02 \x01(\t\"\"\n\x0fGetImageRequest\x12\x0f\n\x07imageID\x18\x01 \x01(\x03\"w\n\x10GetImageResponse\x12\x0e\n\x06\x61uthor\x18\x01 \x01(\t\x12\r\n\x05title\x18\x02 \x01(\t\x12\x13\n\x0b\x64\x65scription\x18\x03 \x01(\t\x12\x0c\n\x04Tags\x18\x04 \x03(\t\x12\x12\n\nimageChunk\x18\x05 \x01(\x0c\x12\r\n\x05\x65rror\x18\x06 \x01(\t\"2\n\x13GetImageListRequest\x12\x0e\n\x06\x61uthor\x18\x01 \x01(\t\x12\x0b\n\x03tag\x18\x02 \x01(\t\"7\n\x14GetImageListResponse\x12\x10\n\x08imageIDs\x18\x01 \x03(\x03\x12\r\n\x05\x65rror\x18\x02 \x01(\t\"v\n\x12ModifyImageRequest\x12\x0f\n\x07imageID\x18\x01 \x01(\x03\x12\r\n\x05token\x18\x02 \x01(\t\x12\x0e\n\x06\x61uthor\x18\x03 \x01(\t\x12\r\n\x05title\x18\x04 \x01(\t\x12\x13\n\x0b\x64\x65scription\x18\x05 \x01(\t\x12\x0c\n\x04Tags\x18\x06 \x03(\t\"$\n\x13ModifyImageResponse\x12\r\n\x05\x65rror\x18\x01 \x01(\t\"4\n\x12\x44\x65leteImageRequest\x12\x0f\n\x07imageID\x18\x01 \x01(\x03\x12\r\n\x05token\x18\x02 \x01(\t\"$\n\x13\x44\x65leteImageResponse\x12\r\n\x05\x65rror\x18\x01 \x01(\t\"\x0e\n\x0c\x45mptyMessageb\x06proto3"

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

UploadImageRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("UploadImageRequest").msgclass
UploadImageResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("UploadImageResponse").msgclass
GetImageRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("GetImageRequest").msgclass
GetImageResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("GetImageResponse").msgclass
GetImageListRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("GetImageListRequest").msgclass
GetImageListResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("GetImageListResponse").msgclass
ModifyImageRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("ModifyImageRequest").msgclass
ModifyImageResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("ModifyImageResponse").msgclass
DeleteImageRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("DeleteImageRequest").msgclass
DeleteImageResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("DeleteImageResponse").msgclass
EmptyMessage = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("EmptyMessage").msgclass

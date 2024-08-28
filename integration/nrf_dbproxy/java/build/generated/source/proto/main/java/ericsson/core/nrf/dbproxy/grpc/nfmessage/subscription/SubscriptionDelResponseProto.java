// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: nfmessage/subscription/SubscriptionDelResponse.proto

package ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription;

public final class SubscriptionDelResponseProto {
  private SubscriptionDelResponseProto() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  public interface SubscriptionDelResponseOrBuilder extends
      // @@protoc_insertion_point(interface_extends:grpc.SubscriptionDelResponse)
      com.google.protobuf.MessageOrBuilder {

    /**
     * <code>uint32 code = 1;</code>
     */
    int getCode();
  }
  /**
   * Protobuf type {@code grpc.SubscriptionDelResponse}
   */
  public  static final class SubscriptionDelResponse extends
      com.google.protobuf.GeneratedMessageV3 implements
      // @@protoc_insertion_point(message_implements:grpc.SubscriptionDelResponse)
      SubscriptionDelResponseOrBuilder {
  private static final long serialVersionUID = 0L;
    // Use SubscriptionDelResponse.newBuilder() to construct.
    private SubscriptionDelResponse(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
      super(builder);
    }
    private SubscriptionDelResponse() {
      code_ = 0;
    }

    @java.lang.Override
    public final com.google.protobuf.UnknownFieldSet
    getUnknownFields() {
      return this.unknownFields;
    }
    private SubscriptionDelResponse(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      this();
      if (extensionRegistry == null) {
        throw new java.lang.NullPointerException();
      }
      int mutable_bitField0_ = 0;
      com.google.protobuf.UnknownFieldSet.Builder unknownFields =
          com.google.protobuf.UnknownFieldSet.newBuilder();
      try {
        boolean done = false;
        while (!done) {
          int tag = input.readTag();
          switch (tag) {
            case 0:
              done = true;
              break;
            default: {
              if (!parseUnknownFieldProto3(
                  input, unknownFields, extensionRegistry, tag)) {
                done = true;
              }
              break;
            }
            case 8: {

              code_ = input.readUInt32();
              break;
            }
          }
        }
      } catch (com.google.protobuf.InvalidProtocolBufferException e) {
        throw e.setUnfinishedMessage(this);
      } catch (java.io.IOException e) {
        throw new com.google.protobuf.InvalidProtocolBufferException(
            e).setUnfinishedMessage(this);
      } finally {
        this.unknownFields = unknownFields.build();
        makeExtensionsImmutable();
      }
    }
    public static final com.google.protobuf.Descriptors.Descriptor
        getDescriptor() {
      return ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.internal_static_grpc_SubscriptionDelResponse_descriptor;
    }

    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.internal_static_grpc_SubscriptionDelResponse_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse.class, ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse.Builder.class);
    }

    public static final int CODE_FIELD_NUMBER = 1;
    private int code_;
    /**
     * <code>uint32 code = 1;</code>
     */
    public int getCode() {
      return code_;
    }

    private byte memoizedIsInitialized = -1;
    public final boolean isInitialized() {
      byte isInitialized = memoizedIsInitialized;
      if (isInitialized == 1) return true;
      if (isInitialized == 0) return false;

      memoizedIsInitialized = 1;
      return true;
    }

    public void writeTo(com.google.protobuf.CodedOutputStream output)
                        throws java.io.IOException {
      if (code_ != 0) {
        output.writeUInt32(1, code_);
      }
      unknownFields.writeTo(output);
    }

    public int getSerializedSize() {
      int size = memoizedSize;
      if (size != -1) return size;

      size = 0;
      if (code_ != 0) {
        size += com.google.protobuf.CodedOutputStream
          .computeUInt32Size(1, code_);
      }
      size += unknownFields.getSerializedSize();
      memoizedSize = size;
      return size;
    }

    @java.lang.Override
    public boolean equals(final java.lang.Object obj) {
      if (obj == this) {
       return true;
      }
      if (!(obj instanceof ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse)) {
        return super.equals(obj);
      }
      ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse other = (ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse) obj;

      boolean result = true;
      result = result && (getCode()
          == other.getCode());
      result = result && unknownFields.equals(other.unknownFields);
      return result;
    }

    @java.lang.Override
    public int hashCode() {
      if (memoizedHashCode != 0) {
        return memoizedHashCode;
      }
      int hash = 41;
      hash = (19 * hash) + getDescriptor().hashCode();
      hash = (37 * hash) + CODE_FIELD_NUMBER;
      hash = (53 * hash) + getCode();
      hash = (29 * hash) + unknownFields.hashCode();
      memoizedHashCode = hash;
      return hash;
    }

    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse parseFrom(
        java.nio.ByteBuffer data)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse parseFrom(
        java.nio.ByteBuffer data,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data, extensionRegistry);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse parseFrom(
        com.google.protobuf.ByteString data)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse parseFrom(
        com.google.protobuf.ByteString data,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data, extensionRegistry);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse parseFrom(byte[] data)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse parseFrom(
        byte[] data,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data, extensionRegistry);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse parseFrom(java.io.InputStream input)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse parseFrom(
        java.io.InputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input, extensionRegistry);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse parseDelimitedFrom(java.io.InputStream input)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseDelimitedWithIOException(PARSER, input);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse parseDelimitedFrom(
        java.io.InputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse parseFrom(
        com.google.protobuf.CodedInputStream input)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse parseFrom(
        com.google.protobuf.CodedInputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input, extensionRegistry);
    }

    public Builder newBuilderForType() { return newBuilder(); }
    public static Builder newBuilder() {
      return DEFAULT_INSTANCE.toBuilder();
    }
    public static Builder newBuilder(ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse prototype) {
      return DEFAULT_INSTANCE.toBuilder().mergeFrom(prototype);
    }
    public Builder toBuilder() {
      return this == DEFAULT_INSTANCE
          ? new Builder() : new Builder().mergeFrom(this);
    }

    @java.lang.Override
    protected Builder newBuilderForType(
        com.google.protobuf.GeneratedMessageV3.BuilderParent parent) {
      Builder builder = new Builder(parent);
      return builder;
    }
    /**
     * Protobuf type {@code grpc.SubscriptionDelResponse}
     */
    public static final class Builder extends
        com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
        // @@protoc_insertion_point(builder_implements:grpc.SubscriptionDelResponse)
        ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponseOrBuilder {
      public static final com.google.protobuf.Descriptors.Descriptor
          getDescriptor() {
        return ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.internal_static_grpc_SubscriptionDelResponse_descriptor;
      }

      protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
          internalGetFieldAccessorTable() {
        return ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.internal_static_grpc_SubscriptionDelResponse_fieldAccessorTable
            .ensureFieldAccessorsInitialized(
                ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse.class, ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse.Builder.class);
      }

      // Construct using ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse.newBuilder()
      private Builder() {
        maybeForceBuilderInitialization();
      }

      private Builder(
          com.google.protobuf.GeneratedMessageV3.BuilderParent parent) {
        super(parent);
        maybeForceBuilderInitialization();
      }
      private void maybeForceBuilderInitialization() {
        if (com.google.protobuf.GeneratedMessageV3
                .alwaysUseFieldBuilders) {
        }
      }
      public Builder clear() {
        super.clear();
        code_ = 0;

        return this;
      }

      public com.google.protobuf.Descriptors.Descriptor
          getDescriptorForType() {
        return ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.internal_static_grpc_SubscriptionDelResponse_descriptor;
      }

      public ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse getDefaultInstanceForType() {
        return ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse.getDefaultInstance();
      }

      public ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse build() {
        ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse result = buildPartial();
        if (!result.isInitialized()) {
          throw newUninitializedMessageException(result);
        }
        return result;
      }

      public ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse buildPartial() {
        ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse result = new ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse(this);
        result.code_ = code_;
        onBuilt();
        return result;
      }

      public Builder clone() {
        return (Builder) super.clone();
      }
      public Builder setField(
          com.google.protobuf.Descriptors.FieldDescriptor field,
          java.lang.Object value) {
        return (Builder) super.setField(field, value);
      }
      public Builder clearField(
          com.google.protobuf.Descriptors.FieldDescriptor field) {
        return (Builder) super.clearField(field);
      }
      public Builder clearOneof(
          com.google.protobuf.Descriptors.OneofDescriptor oneof) {
        return (Builder) super.clearOneof(oneof);
      }
      public Builder setRepeatedField(
          com.google.protobuf.Descriptors.FieldDescriptor field,
          int index, java.lang.Object value) {
        return (Builder) super.setRepeatedField(field, index, value);
      }
      public Builder addRepeatedField(
          com.google.protobuf.Descriptors.FieldDescriptor field,
          java.lang.Object value) {
        return (Builder) super.addRepeatedField(field, value);
      }
      public Builder mergeFrom(com.google.protobuf.Message other) {
        if (other instanceof ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse) {
          return mergeFrom((ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse)other);
        } else {
          super.mergeFrom(other);
          return this;
        }
      }

      public Builder mergeFrom(ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse other) {
        if (other == ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse.getDefaultInstance()) return this;
        if (other.getCode() != 0) {
          setCode(other.getCode());
        }
        this.mergeUnknownFields(other.unknownFields);
        onChanged();
        return this;
      }

      public final boolean isInitialized() {
        return true;
      }

      public Builder mergeFrom(
          com.google.protobuf.CodedInputStream input,
          com.google.protobuf.ExtensionRegistryLite extensionRegistry)
          throws java.io.IOException {
        ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse parsedMessage = null;
        try {
          parsedMessage = PARSER.parsePartialFrom(input, extensionRegistry);
        } catch (com.google.protobuf.InvalidProtocolBufferException e) {
          parsedMessage = (ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse) e.getUnfinishedMessage();
          throw e.unwrapIOException();
        } finally {
          if (parsedMessage != null) {
            mergeFrom(parsedMessage);
          }
        }
        return this;
      }

      private int code_ ;
      /**
       * <code>uint32 code = 1;</code>
       */
      public int getCode() {
        return code_;
      }
      /**
       * <code>uint32 code = 1;</code>
       */
      public Builder setCode(int value) {
        
        code_ = value;
        onChanged();
        return this;
      }
      /**
       * <code>uint32 code = 1;</code>
       */
      public Builder clearCode() {
        
        code_ = 0;
        onChanged();
        return this;
      }
      public final Builder setUnknownFields(
          final com.google.protobuf.UnknownFieldSet unknownFields) {
        return super.setUnknownFieldsProto3(unknownFields);
      }

      public final Builder mergeUnknownFields(
          final com.google.protobuf.UnknownFieldSet unknownFields) {
        return super.mergeUnknownFields(unknownFields);
      }


      // @@protoc_insertion_point(builder_scope:grpc.SubscriptionDelResponse)
    }

    // @@protoc_insertion_point(class_scope:grpc.SubscriptionDelResponse)
    private static final ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse DEFAULT_INSTANCE;
    static {
      DEFAULT_INSTANCE = new ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse();
    }

    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse getDefaultInstance() {
      return DEFAULT_INSTANCE;
    }

    private static final com.google.protobuf.Parser<SubscriptionDelResponse>
        PARSER = new com.google.protobuf.AbstractParser<SubscriptionDelResponse>() {
      public SubscriptionDelResponse parsePartialFrom(
          com.google.protobuf.CodedInputStream input,
          com.google.protobuf.ExtensionRegistryLite extensionRegistry)
          throws com.google.protobuf.InvalidProtocolBufferException {
        return new SubscriptionDelResponse(input, extensionRegistry);
      }
    };

    public static com.google.protobuf.Parser<SubscriptionDelResponse> parser() {
      return PARSER;
    }

    @java.lang.Override
    public com.google.protobuf.Parser<SubscriptionDelResponse> getParserForType() {
      return PARSER;
    }

    public ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse getDefaultInstanceForType() {
      return DEFAULT_INSTANCE;
    }

  }

  private static final com.google.protobuf.Descriptors.Descriptor
    internal_static_grpc_SubscriptionDelResponse_descriptor;
  private static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_grpc_SubscriptionDelResponse_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n4nfmessage/subscription/SubscriptionDel" +
      "Response.proto\022\004grpc\"\'\n\027SubscriptionDelR" +
      "esponse\022\014\n\004code\030\001 \001(\rBy\n5ericsson.core.n" +
      "rf.dbproxy.grpc.nfmessage.subscriptionB\034" +
      "SubscriptionDelResponseProtoZ\"com/dbprox" +
      "y/nfmessage/subscriptionb\006proto3"
    };
    com.google.protobuf.Descriptors.FileDescriptor.InternalDescriptorAssigner assigner =
        new com.google.protobuf.Descriptors.FileDescriptor.    InternalDescriptorAssigner() {
          public com.google.protobuf.ExtensionRegistry assignDescriptors(
              com.google.protobuf.Descriptors.FileDescriptor root) {
            descriptor = root;
            return null;
          }
        };
    com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
        }, assigner);
    internal_static_grpc_SubscriptionDelResponse_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_grpc_SubscriptionDelResponse_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_grpc_SubscriptionDelResponse_descriptor,
        new java.lang.String[] { "Code", });
  }

  // @@protoc_insertion_point(outer_class_scope)
}

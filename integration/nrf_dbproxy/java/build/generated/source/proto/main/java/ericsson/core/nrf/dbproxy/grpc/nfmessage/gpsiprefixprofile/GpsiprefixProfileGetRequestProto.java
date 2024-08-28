// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: nfmessage/gpsiprefixprofile/GpsiprefixProfileGetRequest.proto

package ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile;

public final class GpsiprefixProfileGetRequestProto {
  private GpsiprefixProfileGetRequestProto() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  public interface GpsiprefixProfileGetRequestOrBuilder extends
      // @@protoc_insertion_point(interface_extends:grpc.GpsiprefixProfileGetRequest)
      com.google.protobuf.MessageOrBuilder {

    /**
     * <code>uint64 search_gpsi = 1;</code>
     */
    long getSearchGpsi();
  }
  /**
   * Protobuf type {@code grpc.GpsiprefixProfileGetRequest}
   */
  public  static final class GpsiprefixProfileGetRequest extends
      com.google.protobuf.GeneratedMessageV3 implements
      // @@protoc_insertion_point(message_implements:grpc.GpsiprefixProfileGetRequest)
      GpsiprefixProfileGetRequestOrBuilder {
  private static final long serialVersionUID = 0L;
    // Use GpsiprefixProfileGetRequest.newBuilder() to construct.
    private GpsiprefixProfileGetRequest(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
      super(builder);
    }
    private GpsiprefixProfileGetRequest() {
      searchGpsi_ = 0L;
    }

    @java.lang.Override
    public final com.google.protobuf.UnknownFieldSet
    getUnknownFields() {
      return this.unknownFields;
    }
    private GpsiprefixProfileGetRequest(
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

              searchGpsi_ = input.readUInt64();
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
      return ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.internal_static_grpc_GpsiprefixProfileGetRequest_descriptor;
    }

    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.internal_static_grpc_GpsiprefixProfileGetRequest_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest.class, ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest.Builder.class);
    }

    public static final int SEARCH_GPSI_FIELD_NUMBER = 1;
    private long searchGpsi_;
    /**
     * <code>uint64 search_gpsi = 1;</code>
     */
    public long getSearchGpsi() {
      return searchGpsi_;
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
      if (searchGpsi_ != 0L) {
        output.writeUInt64(1, searchGpsi_);
      }
      unknownFields.writeTo(output);
    }

    public int getSerializedSize() {
      int size = memoizedSize;
      if (size != -1) return size;

      size = 0;
      if (searchGpsi_ != 0L) {
        size += com.google.protobuf.CodedOutputStream
          .computeUInt64Size(1, searchGpsi_);
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
      if (!(obj instanceof ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest)) {
        return super.equals(obj);
      }
      ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest other = (ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest) obj;

      boolean result = true;
      result = result && (getSearchGpsi()
          == other.getSearchGpsi());
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
      hash = (37 * hash) + SEARCH_GPSI_FIELD_NUMBER;
      hash = (53 * hash) + com.google.protobuf.Internal.hashLong(
          getSearchGpsi());
      hash = (29 * hash) + unknownFields.hashCode();
      memoizedHashCode = hash;
      return hash;
    }

    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest parseFrom(
        java.nio.ByteBuffer data)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest parseFrom(
        java.nio.ByteBuffer data,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data, extensionRegistry);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest parseFrom(
        com.google.protobuf.ByteString data)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest parseFrom(
        com.google.protobuf.ByteString data,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data, extensionRegistry);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest parseFrom(byte[] data)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest parseFrom(
        byte[] data,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data, extensionRegistry);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest parseFrom(java.io.InputStream input)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest parseFrom(
        java.io.InputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input, extensionRegistry);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest parseDelimitedFrom(java.io.InputStream input)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseDelimitedWithIOException(PARSER, input);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest parseDelimitedFrom(
        java.io.InputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest parseFrom(
        com.google.protobuf.CodedInputStream input)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest parseFrom(
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
    public static Builder newBuilder(ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest prototype) {
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
     * Protobuf type {@code grpc.GpsiprefixProfileGetRequest}
     */
    public static final class Builder extends
        com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
        // @@protoc_insertion_point(builder_implements:grpc.GpsiprefixProfileGetRequest)
        ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequestOrBuilder {
      public static final com.google.protobuf.Descriptors.Descriptor
          getDescriptor() {
        return ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.internal_static_grpc_GpsiprefixProfileGetRequest_descriptor;
      }

      protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
          internalGetFieldAccessorTable() {
        return ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.internal_static_grpc_GpsiprefixProfileGetRequest_fieldAccessorTable
            .ensureFieldAccessorsInitialized(
                ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest.class, ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest.Builder.class);
      }

      // Construct using ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest.newBuilder()
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
        searchGpsi_ = 0L;

        return this;
      }

      public com.google.protobuf.Descriptors.Descriptor
          getDescriptorForType() {
        return ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.internal_static_grpc_GpsiprefixProfileGetRequest_descriptor;
      }

      public ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest getDefaultInstanceForType() {
        return ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest.getDefaultInstance();
      }

      public ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest build() {
        ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest result = buildPartial();
        if (!result.isInitialized()) {
          throw newUninitializedMessageException(result);
        }
        return result;
      }

      public ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest buildPartial() {
        ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest result = new ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest(this);
        result.searchGpsi_ = searchGpsi_;
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
        if (other instanceof ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest) {
          return mergeFrom((ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest)other);
        } else {
          super.mergeFrom(other);
          return this;
        }
      }

      public Builder mergeFrom(ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest other) {
        if (other == ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest.getDefaultInstance()) return this;
        if (other.getSearchGpsi() != 0L) {
          setSearchGpsi(other.getSearchGpsi());
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
        ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest parsedMessage = null;
        try {
          parsedMessage = PARSER.parsePartialFrom(input, extensionRegistry);
        } catch (com.google.protobuf.InvalidProtocolBufferException e) {
          parsedMessage = (ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest) e.getUnfinishedMessage();
          throw e.unwrapIOException();
        } finally {
          if (parsedMessage != null) {
            mergeFrom(parsedMessage);
          }
        }
        return this;
      }

      private long searchGpsi_ ;
      /**
       * <code>uint64 search_gpsi = 1;</code>
       */
      public long getSearchGpsi() {
        return searchGpsi_;
      }
      /**
       * <code>uint64 search_gpsi = 1;</code>
       */
      public Builder setSearchGpsi(long value) {
        
        searchGpsi_ = value;
        onChanged();
        return this;
      }
      /**
       * <code>uint64 search_gpsi = 1;</code>
       */
      public Builder clearSearchGpsi() {
        
        searchGpsi_ = 0L;
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


      // @@protoc_insertion_point(builder_scope:grpc.GpsiprefixProfileGetRequest)
    }

    // @@protoc_insertion_point(class_scope:grpc.GpsiprefixProfileGetRequest)
    private static final ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest DEFAULT_INSTANCE;
    static {
      DEFAULT_INSTANCE = new ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest();
    }

    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest getDefaultInstance() {
      return DEFAULT_INSTANCE;
    }

    private static final com.google.protobuf.Parser<GpsiprefixProfileGetRequest>
        PARSER = new com.google.protobuf.AbstractParser<GpsiprefixProfileGetRequest>() {
      public GpsiprefixProfileGetRequest parsePartialFrom(
          com.google.protobuf.CodedInputStream input,
          com.google.protobuf.ExtensionRegistryLite extensionRegistry)
          throws com.google.protobuf.InvalidProtocolBufferException {
        return new GpsiprefixProfileGetRequest(input, extensionRegistry);
      }
    };

    public static com.google.protobuf.Parser<GpsiprefixProfileGetRequest> parser() {
      return PARSER;
    }

    @java.lang.Override
    public com.google.protobuf.Parser<GpsiprefixProfileGetRequest> getParserForType() {
      return PARSER;
    }

    public ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest getDefaultInstanceForType() {
      return DEFAULT_INSTANCE;
    }

  }

  private static final com.google.protobuf.Descriptors.Descriptor
    internal_static_grpc_GpsiprefixProfileGetRequest_descriptor;
  private static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_grpc_GpsiprefixProfileGetRequest_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n=nfmessage/gpsiprefixprofile/Gpsiprefix" +
      "ProfileGetRequest.proto\022\004grpc\"2\n\033Gpsipre" +
      "fixProfileGetRequest\022\023\n\013search_gpsi\030\001 \001(" +
      "\004B\207\001\n:ericsson.core.nrf.dbproxy.grpc.nfm" +
      "essage.gpsiprefixprofileB GpsiprefixProf" +
      "ileGetRequestProtoZ\'com/dbproxy/nfmessag" +
      "e/gpsiprefixprofileb\006proto3"
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
    internal_static_grpc_GpsiprefixProfileGetRequest_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_grpc_GpsiprefixProfileGetRequest_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_grpc_GpsiprefixProfileGetRequest_descriptor,
        new java.lang.String[] { "SearchGpsi", });
  }

  // @@protoc_insertion_point(outer_class_scope)
}

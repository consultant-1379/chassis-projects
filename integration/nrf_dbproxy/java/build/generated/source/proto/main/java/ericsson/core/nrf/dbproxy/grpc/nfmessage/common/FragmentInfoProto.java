// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: nfmessage/common/FragmentInfo.proto

package ericsson.core.nrf.dbproxy.grpc.nfmessage.common;

public final class FragmentInfoProto {
  private FragmentInfoProto() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  public interface FragmentInfoOrBuilder extends
      // @@protoc_insertion_point(interface_extends:grpc.FragmentInfo)
      com.google.protobuf.MessageOrBuilder {

    /**
     * <code>uint32 total_number = 1;</code>
     */
    int getTotalNumber();

    /**
     * <code>uint32 transmitted_number = 2;</code>
     */
    int getTransmittedNumber();

    /**
     * <code>string fragment_session_id = 3;</code>
     */
    java.lang.String getFragmentSessionId();
    /**
     * <code>string fragment_session_id = 3;</code>
     */
    com.google.protobuf.ByteString
        getFragmentSessionIdBytes();
  }
  /**
   * Protobuf type {@code grpc.FragmentInfo}
   */
  public  static final class FragmentInfo extends
      com.google.protobuf.GeneratedMessageV3 implements
      // @@protoc_insertion_point(message_implements:grpc.FragmentInfo)
      FragmentInfoOrBuilder {
  private static final long serialVersionUID = 0L;
    // Use FragmentInfo.newBuilder() to construct.
    private FragmentInfo(com.google.protobuf.GeneratedMessageV3.Builder<?> builder) {
      super(builder);
    }
    private FragmentInfo() {
      totalNumber_ = 0;
      transmittedNumber_ = 0;
      fragmentSessionId_ = "";
    }

    @java.lang.Override
    public final com.google.protobuf.UnknownFieldSet
    getUnknownFields() {
      return this.unknownFields;
    }
    private FragmentInfo(
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

              totalNumber_ = input.readUInt32();
              break;
            }
            case 16: {

              transmittedNumber_ = input.readUInt32();
              break;
            }
            case 26: {
              java.lang.String s = input.readStringRequireUtf8();

              fragmentSessionId_ = s;
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
      return ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.internal_static_grpc_FragmentInfo_descriptor;
    }

    protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
        internalGetFieldAccessorTable() {
      return ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.internal_static_grpc_FragmentInfo_fieldAccessorTable
          .ensureFieldAccessorsInitialized(
              ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo.class, ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo.Builder.class);
    }

    public static final int TOTAL_NUMBER_FIELD_NUMBER = 1;
    private int totalNumber_;
    /**
     * <code>uint32 total_number = 1;</code>
     */
    public int getTotalNumber() {
      return totalNumber_;
    }

    public static final int TRANSMITTED_NUMBER_FIELD_NUMBER = 2;
    private int transmittedNumber_;
    /**
     * <code>uint32 transmitted_number = 2;</code>
     */
    public int getTransmittedNumber() {
      return transmittedNumber_;
    }

    public static final int FRAGMENT_SESSION_ID_FIELD_NUMBER = 3;
    private volatile java.lang.Object fragmentSessionId_;
    /**
     * <code>string fragment_session_id = 3;</code>
     */
    public java.lang.String getFragmentSessionId() {
      java.lang.Object ref = fragmentSessionId_;
      if (ref instanceof java.lang.String) {
        return (java.lang.String) ref;
      } else {
        com.google.protobuf.ByteString bs = 
            (com.google.protobuf.ByteString) ref;
        java.lang.String s = bs.toStringUtf8();
        fragmentSessionId_ = s;
        return s;
      }
    }
    /**
     * <code>string fragment_session_id = 3;</code>
     */
    public com.google.protobuf.ByteString
        getFragmentSessionIdBytes() {
      java.lang.Object ref = fragmentSessionId_;
      if (ref instanceof java.lang.String) {
        com.google.protobuf.ByteString b = 
            com.google.protobuf.ByteString.copyFromUtf8(
                (java.lang.String) ref);
        fragmentSessionId_ = b;
        return b;
      } else {
        return (com.google.protobuf.ByteString) ref;
      }
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
      if (totalNumber_ != 0) {
        output.writeUInt32(1, totalNumber_);
      }
      if (transmittedNumber_ != 0) {
        output.writeUInt32(2, transmittedNumber_);
      }
      if (!getFragmentSessionIdBytes().isEmpty()) {
        com.google.protobuf.GeneratedMessageV3.writeString(output, 3, fragmentSessionId_);
      }
      unknownFields.writeTo(output);
    }

    public int getSerializedSize() {
      int size = memoizedSize;
      if (size != -1) return size;

      size = 0;
      if (totalNumber_ != 0) {
        size += com.google.protobuf.CodedOutputStream
          .computeUInt32Size(1, totalNumber_);
      }
      if (transmittedNumber_ != 0) {
        size += com.google.protobuf.CodedOutputStream
          .computeUInt32Size(2, transmittedNumber_);
      }
      if (!getFragmentSessionIdBytes().isEmpty()) {
        size += com.google.protobuf.GeneratedMessageV3.computeStringSize(3, fragmentSessionId_);
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
      if (!(obj instanceof ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo)) {
        return super.equals(obj);
      }
      ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo other = (ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo) obj;

      boolean result = true;
      result = result && (getTotalNumber()
          == other.getTotalNumber());
      result = result && (getTransmittedNumber()
          == other.getTransmittedNumber());
      result = result && getFragmentSessionId()
          .equals(other.getFragmentSessionId());
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
      hash = (37 * hash) + TOTAL_NUMBER_FIELD_NUMBER;
      hash = (53 * hash) + getTotalNumber();
      hash = (37 * hash) + TRANSMITTED_NUMBER_FIELD_NUMBER;
      hash = (53 * hash) + getTransmittedNumber();
      hash = (37 * hash) + FRAGMENT_SESSION_ID_FIELD_NUMBER;
      hash = (53 * hash) + getFragmentSessionId().hashCode();
      hash = (29 * hash) + unknownFields.hashCode();
      memoizedHashCode = hash;
      return hash;
    }

    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo parseFrom(
        java.nio.ByteBuffer data)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo parseFrom(
        java.nio.ByteBuffer data,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data, extensionRegistry);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo parseFrom(
        com.google.protobuf.ByteString data)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo parseFrom(
        com.google.protobuf.ByteString data,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data, extensionRegistry);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo parseFrom(byte[] data)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo parseFrom(
        byte[] data,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws com.google.protobuf.InvalidProtocolBufferException {
      return PARSER.parseFrom(data, extensionRegistry);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo parseFrom(java.io.InputStream input)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo parseFrom(
        java.io.InputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input, extensionRegistry);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo parseDelimitedFrom(java.io.InputStream input)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseDelimitedWithIOException(PARSER, input);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo parseDelimitedFrom(
        java.io.InputStream input,
        com.google.protobuf.ExtensionRegistryLite extensionRegistry)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseDelimitedWithIOException(PARSER, input, extensionRegistry);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo parseFrom(
        com.google.protobuf.CodedInputStream input)
        throws java.io.IOException {
      return com.google.protobuf.GeneratedMessageV3
          .parseWithIOException(PARSER, input);
    }
    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo parseFrom(
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
    public static Builder newBuilder(ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo prototype) {
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
     * Protobuf type {@code grpc.FragmentInfo}
     */
    public static final class Builder extends
        com.google.protobuf.GeneratedMessageV3.Builder<Builder> implements
        // @@protoc_insertion_point(builder_implements:grpc.FragmentInfo)
        ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfoOrBuilder {
      public static final com.google.protobuf.Descriptors.Descriptor
          getDescriptor() {
        return ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.internal_static_grpc_FragmentInfo_descriptor;
      }

      protected com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
          internalGetFieldAccessorTable() {
        return ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.internal_static_grpc_FragmentInfo_fieldAccessorTable
            .ensureFieldAccessorsInitialized(
                ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo.class, ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo.Builder.class);
      }

      // Construct using ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo.newBuilder()
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
        totalNumber_ = 0;

        transmittedNumber_ = 0;

        fragmentSessionId_ = "";

        return this;
      }

      public com.google.protobuf.Descriptors.Descriptor
          getDescriptorForType() {
        return ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.internal_static_grpc_FragmentInfo_descriptor;
      }

      public ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo getDefaultInstanceForType() {
        return ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo.getDefaultInstance();
      }

      public ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo build() {
        ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo result = buildPartial();
        if (!result.isInitialized()) {
          throw newUninitializedMessageException(result);
        }
        return result;
      }

      public ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo buildPartial() {
        ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo result = new ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo(this);
        result.totalNumber_ = totalNumber_;
        result.transmittedNumber_ = transmittedNumber_;
        result.fragmentSessionId_ = fragmentSessionId_;
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
        if (other instanceof ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo) {
          return mergeFrom((ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo)other);
        } else {
          super.mergeFrom(other);
          return this;
        }
      }

      public Builder mergeFrom(ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo other) {
        if (other == ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo.getDefaultInstance()) return this;
        if (other.getTotalNumber() != 0) {
          setTotalNumber(other.getTotalNumber());
        }
        if (other.getTransmittedNumber() != 0) {
          setTransmittedNumber(other.getTransmittedNumber());
        }
        if (!other.getFragmentSessionId().isEmpty()) {
          fragmentSessionId_ = other.fragmentSessionId_;
          onChanged();
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
        ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo parsedMessage = null;
        try {
          parsedMessage = PARSER.parsePartialFrom(input, extensionRegistry);
        } catch (com.google.protobuf.InvalidProtocolBufferException e) {
          parsedMessage = (ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo) e.getUnfinishedMessage();
          throw e.unwrapIOException();
        } finally {
          if (parsedMessage != null) {
            mergeFrom(parsedMessage);
          }
        }
        return this;
      }

      private int totalNumber_ ;
      /**
       * <code>uint32 total_number = 1;</code>
       */
      public int getTotalNumber() {
        return totalNumber_;
      }
      /**
       * <code>uint32 total_number = 1;</code>
       */
      public Builder setTotalNumber(int value) {
        
        totalNumber_ = value;
        onChanged();
        return this;
      }
      /**
       * <code>uint32 total_number = 1;</code>
       */
      public Builder clearTotalNumber() {
        
        totalNumber_ = 0;
        onChanged();
        return this;
      }

      private int transmittedNumber_ ;
      /**
       * <code>uint32 transmitted_number = 2;</code>
       */
      public int getTransmittedNumber() {
        return transmittedNumber_;
      }
      /**
       * <code>uint32 transmitted_number = 2;</code>
       */
      public Builder setTransmittedNumber(int value) {
        
        transmittedNumber_ = value;
        onChanged();
        return this;
      }
      /**
       * <code>uint32 transmitted_number = 2;</code>
       */
      public Builder clearTransmittedNumber() {
        
        transmittedNumber_ = 0;
        onChanged();
        return this;
      }

      private java.lang.Object fragmentSessionId_ = "";
      /**
       * <code>string fragment_session_id = 3;</code>
       */
      public java.lang.String getFragmentSessionId() {
        java.lang.Object ref = fragmentSessionId_;
        if (!(ref instanceof java.lang.String)) {
          com.google.protobuf.ByteString bs =
              (com.google.protobuf.ByteString) ref;
          java.lang.String s = bs.toStringUtf8();
          fragmentSessionId_ = s;
          return s;
        } else {
          return (java.lang.String) ref;
        }
      }
      /**
       * <code>string fragment_session_id = 3;</code>
       */
      public com.google.protobuf.ByteString
          getFragmentSessionIdBytes() {
        java.lang.Object ref = fragmentSessionId_;
        if (ref instanceof String) {
          com.google.protobuf.ByteString b = 
              com.google.protobuf.ByteString.copyFromUtf8(
                  (java.lang.String) ref);
          fragmentSessionId_ = b;
          return b;
        } else {
          return (com.google.protobuf.ByteString) ref;
        }
      }
      /**
       * <code>string fragment_session_id = 3;</code>
       */
      public Builder setFragmentSessionId(
          java.lang.String value) {
        if (value == null) {
    throw new NullPointerException();
  }
  
        fragmentSessionId_ = value;
        onChanged();
        return this;
      }
      /**
       * <code>string fragment_session_id = 3;</code>
       */
      public Builder clearFragmentSessionId() {
        
        fragmentSessionId_ = getDefaultInstance().getFragmentSessionId();
        onChanged();
        return this;
      }
      /**
       * <code>string fragment_session_id = 3;</code>
       */
      public Builder setFragmentSessionIdBytes(
          com.google.protobuf.ByteString value) {
        if (value == null) {
    throw new NullPointerException();
  }
  checkByteStringIsUtf8(value);
        
        fragmentSessionId_ = value;
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


      // @@protoc_insertion_point(builder_scope:grpc.FragmentInfo)
    }

    // @@protoc_insertion_point(class_scope:grpc.FragmentInfo)
    private static final ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo DEFAULT_INSTANCE;
    static {
      DEFAULT_INSTANCE = new ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo();
    }

    public static ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo getDefaultInstance() {
      return DEFAULT_INSTANCE;
    }

    private static final com.google.protobuf.Parser<FragmentInfo>
        PARSER = new com.google.protobuf.AbstractParser<FragmentInfo>() {
      public FragmentInfo parsePartialFrom(
          com.google.protobuf.CodedInputStream input,
          com.google.protobuf.ExtensionRegistryLite extensionRegistry)
          throws com.google.protobuf.InvalidProtocolBufferException {
        return new FragmentInfo(input, extensionRegistry);
      }
    };

    public static com.google.protobuf.Parser<FragmentInfo> parser() {
      return PARSER;
    }

    @java.lang.Override
    public com.google.protobuf.Parser<FragmentInfo> getParserForType() {
      return PARSER;
    }

    public ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo getDefaultInstanceForType() {
      return DEFAULT_INSTANCE;
    }

  }

  private static final com.google.protobuf.Descriptors.Descriptor
    internal_static_grpc_FragmentInfo_descriptor;
  private static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_grpc_FragmentInfo_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n#nfmessage/common/FragmentInfo.proto\022\004g" +
      "rpc\"]\n\014FragmentInfo\022\024\n\014total_number\030\001 \001(" +
      "\r\022\032\n\022transmitted_number\030\002 \001(\r\022\033\n\023fragmen" +
      "t_session_id\030\003 \001(\tBb\n/ericsson.core.nrf." +
      "dbproxy.grpc.nfmessage.commonB\021FragmentI" +
      "nfoProtoZ\034com/dbproxy/nfmessage/commonb\006" +
      "proto3"
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
    internal_static_grpc_FragmentInfo_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_grpc_FragmentInfo_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_grpc_FragmentInfo_descriptor,
        new java.lang.String[] { "TotalNumber", "TransmittedNumber", "FragmentSessionId", });
  }

  // @@protoc_insertion_point(outer_class_scope)
}

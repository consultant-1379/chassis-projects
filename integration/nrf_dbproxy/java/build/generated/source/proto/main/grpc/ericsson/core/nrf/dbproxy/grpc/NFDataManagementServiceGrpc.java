package ericsson.core.nrf.dbproxy.grpc;

import static io.grpc.MethodDescriptor.generateFullMethodName;
import static io.grpc.stub.ClientCalls.asyncBidiStreamingCall;
import static io.grpc.stub.ClientCalls.asyncClientStreamingCall;
import static io.grpc.stub.ClientCalls.asyncServerStreamingCall;
import static io.grpc.stub.ClientCalls.asyncUnaryCall;
import static io.grpc.stub.ClientCalls.blockingServerStreamingCall;
import static io.grpc.stub.ClientCalls.blockingUnaryCall;
import static io.grpc.stub.ClientCalls.futureUnaryCall;
import static io.grpc.stub.ServerCalls.asyncBidiStreamingCall;
import static io.grpc.stub.ServerCalls.asyncClientStreamingCall;
import static io.grpc.stub.ServerCalls.asyncServerStreamingCall;
import static io.grpc.stub.ServerCalls.asyncUnaryCall;
import static io.grpc.stub.ServerCalls.asyncUnimplementedStreamingCall;
import static io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.11.0)",
    comments = "Source: NFDataManagementService.proto")
public final class NFDataManagementServiceGrpc {

  private NFDataManagementServiceGrpc() {}

  public static final String SERVICE_NAME = "grpc.NFDataManagementService";

  // Static method descriptors that strictly reflect the proto.
  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  @java.lang.Deprecated // Use {@link #getExecuteMethod()} instead. 
  public static final io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage> METHOD_EXECUTE = getExecuteMethodHelper();

  private static volatile io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage> getExecuteMethod;

  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  public static io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage> getExecuteMethod() {
    return getExecuteMethodHelper();
  }

  private static io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage> getExecuteMethodHelper() {
    io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage, ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage> getExecuteMethod;
    if ((getExecuteMethod = NFDataManagementServiceGrpc.getExecuteMethod) == null) {
      synchronized (NFDataManagementServiceGrpc.class) {
        if ((getExecuteMethod = NFDataManagementServiceGrpc.getExecuteMethod) == null) {
          NFDataManagementServiceGrpc.getExecuteMethod = getExecuteMethod = 
              io.grpc.MethodDescriptor.<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage, ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "grpc.NFDataManagementService", "execute"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage.getDefaultInstance()))
                  .setSchemaDescriptor(new NFDataManagementServiceMethodDescriptorSupplier("execute"))
                  .build();
          }
        }
     }
     return getExecuteMethod;
  }
  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  @java.lang.Deprecated // Use {@link #getInsertMethod()} instead. 
  public static final io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertResponse> METHOD_INSERT = getInsertMethodHelper();

  private static volatile io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertResponse> getInsertMethod;

  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  public static io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertResponse> getInsertMethod() {
    return getInsertMethodHelper();
  }

  private static io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertResponse> getInsertMethodHelper() {
    io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertRequest, ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertResponse> getInsertMethod;
    if ((getInsertMethod = NFDataManagementServiceGrpc.getInsertMethod) == null) {
      synchronized (NFDataManagementServiceGrpc.class) {
        if ((getInsertMethod = NFDataManagementServiceGrpc.getInsertMethod) == null) {
          NFDataManagementServiceGrpc.getInsertMethod = getInsertMethod = 
              io.grpc.MethodDescriptor.<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertRequest, ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "grpc.NFDataManagementService", "insert"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new NFDataManagementServiceMethodDescriptorSupplier("insert"))
                  .build();
          }
        }
     }
     return getInsertMethod;
  }
  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  @java.lang.Deprecated // Use {@link #getRemoveMethod()} instead. 
  public static final io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveResponse> METHOD_REMOVE = getRemoveMethodHelper();

  private static volatile io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveResponse> getRemoveMethod;

  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  public static io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveResponse> getRemoveMethod() {
    return getRemoveMethodHelper();
  }

  private static io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveResponse> getRemoveMethodHelper() {
    io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveRequest, ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveResponse> getRemoveMethod;
    if ((getRemoveMethod = NFDataManagementServiceGrpc.getRemoveMethod) == null) {
      synchronized (NFDataManagementServiceGrpc.class) {
        if ((getRemoveMethod = NFDataManagementServiceGrpc.getRemoveMethod) == null) {
          NFDataManagementServiceGrpc.getRemoveMethod = getRemoveMethod = 
              io.grpc.MethodDescriptor.<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveRequest, ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "grpc.NFDataManagementService", "remove"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new NFDataManagementServiceMethodDescriptorSupplier("remove"))
                  .build();
          }
        }
     }
     return getRemoveMethod;
  }
  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  @java.lang.Deprecated // Use {@link #getQueryByKeyMethod()} instead. 
  public static final io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse> METHOD_QUERY_BY_KEY = getQueryByKeyMethodHelper();

  private static volatile io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse> getQueryByKeyMethod;

  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  public static io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse> getQueryByKeyMethod() {
    return getQueryByKeyMethodHelper();
  }

  private static io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse> getQueryByKeyMethodHelper() {
    io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest, ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse> getQueryByKeyMethod;
    if ((getQueryByKeyMethod = NFDataManagementServiceGrpc.getQueryByKeyMethod) == null) {
      synchronized (NFDataManagementServiceGrpc.class) {
        if ((getQueryByKeyMethod = NFDataManagementServiceGrpc.getQueryByKeyMethod) == null) {
          NFDataManagementServiceGrpc.getQueryByKeyMethod = getQueryByKeyMethod = 
              io.grpc.MethodDescriptor.<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest, ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
              .setFullMethodName(generateFullMethodName(
                  "grpc.NFDataManagementService", "queryByKey"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new NFDataManagementServiceMethodDescriptorSupplier("queryByKey"))
                  .build();
          }
        }
     }
     return getQueryByKeyMethod;
  }
  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  @java.lang.Deprecated // Use {@link #getQueryByFilterMethod()} instead. 
  public static final io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse> METHOD_QUERY_BY_FILTER = getQueryByFilterMethodHelper();

  private static volatile io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse> getQueryByFilterMethod;

  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  public static io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse> getQueryByFilterMethod() {
    return getQueryByFilterMethodHelper();
  }

  private static io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse> getQueryByFilterMethodHelper() {
    io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest, ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse> getQueryByFilterMethod;
    if ((getQueryByFilterMethod = NFDataManagementServiceGrpc.getQueryByFilterMethod) == null) {
      synchronized (NFDataManagementServiceGrpc.class) {
        if ((getQueryByFilterMethod = NFDataManagementServiceGrpc.getQueryByFilterMethod) == null) {
          NFDataManagementServiceGrpc.getQueryByFilterMethod = getQueryByFilterMethod = 
              io.grpc.MethodDescriptor.<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest, ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
              .setFullMethodName(generateFullMethodName(
                  "grpc.NFDataManagementService", "queryByFilter"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new NFDataManagementServiceMethodDescriptorSupplier("queryByFilter"))
                  .build();
          }
        }
     }
     return getQueryByFilterMethod;
  }
  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  @java.lang.Deprecated // Use {@link #getTransferParameterMethod()} instead. 
  public static final io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaResponse> METHOD_TRANSFER_PARAMETER = getTransferParameterMethodHelper();

  private static volatile io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaResponse> getTransferParameterMethod;

  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  public static io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaResponse> getTransferParameterMethod() {
    return getTransferParameterMethodHelper();
  }

  private static io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaResponse> getTransferParameterMethodHelper() {
    io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaRequest, ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaResponse> getTransferParameterMethod;
    if ((getTransferParameterMethod = NFDataManagementServiceGrpc.getTransferParameterMethod) == null) {
      synchronized (NFDataManagementServiceGrpc.class) {
        if ((getTransferParameterMethod = NFDataManagementServiceGrpc.getTransferParameterMethod) == null) {
          NFDataManagementServiceGrpc.getTransferParameterMethod = getTransferParameterMethod = 
              io.grpc.MethodDescriptor.<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaRequest, ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "grpc.NFDataManagementService", "transferParameter"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new NFDataManagementServiceMethodDescriptorSupplier("transferParameter"))
                  .build();
          }
        }
     }
     return getTransferParameterMethod;
  }
  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  @java.lang.Deprecated // Use {@link #getPatchNrfProfileMethod()} instead. 
  public static final io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchResponse> METHOD_PATCH_NRF_PROFILE = getPatchNrfProfileMethodHelper();

  private static volatile io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchResponse> getPatchNrfProfileMethod;

  @io.grpc.ExperimentalApi("https://github.com/grpc/grpc-java/issues/1901")
  public static io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchResponse> getPatchNrfProfileMethod() {
    return getPatchNrfProfileMethodHelper();
  }

  private static io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchRequest,
      ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchResponse> getPatchNrfProfileMethodHelper() {
    io.grpc.MethodDescriptor<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchRequest, ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchResponse> getPatchNrfProfileMethod;
    if ((getPatchNrfProfileMethod = NFDataManagementServiceGrpc.getPatchNrfProfileMethod) == null) {
      synchronized (NFDataManagementServiceGrpc.class) {
        if ((getPatchNrfProfileMethod = NFDataManagementServiceGrpc.getPatchNrfProfileMethod) == null) {
          NFDataManagementServiceGrpc.getPatchNrfProfileMethod = getPatchNrfProfileMethod = 
              io.grpc.MethodDescriptor.<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchRequest, ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "grpc.NFDataManagementService", "patchNrfProfile"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new NFDataManagementServiceMethodDescriptorSupplier("patchNrfProfile"))
                  .build();
          }
        }
     }
     return getPatchNrfProfileMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static NFDataManagementServiceStub newStub(io.grpc.Channel channel) {
    return new NFDataManagementServiceStub(channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static NFDataManagementServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    return new NFDataManagementServiceBlockingStub(channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static NFDataManagementServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    return new NFDataManagementServiceFutureStub(channel);
  }

  /**
   */
  public static abstract class NFDataManagementServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void execute(ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage request,
        io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage> responseObserver) {
      asyncUnimplementedUnaryCall(getExecuteMethodHelper(), responseObserver);
    }

    /**
     */
    public void insert(ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertRequest request,
        io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getInsertMethodHelper(), responseObserver);
    }

    /**
     */
    public void remove(ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveRequest request,
        io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getRemoveMethodHelper(), responseObserver);
    }

    /**
     */
    public void queryByKey(ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest request,
        io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getQueryByKeyMethodHelper(), responseObserver);
    }

    /**
     */
    public void queryByFilter(ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest request,
        io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getQueryByFilterMethodHelper(), responseObserver);
    }

    /**
     */
    public void transferParameter(ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaRequest request,
        io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getTransferParameterMethodHelper(), responseObserver);
    }

    /**
     */
    public void patchNrfProfile(ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchRequest request,
        io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getPatchNrfProfileMethodHelper(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getExecuteMethodHelper(),
            asyncUnaryCall(
              new MethodHandlers<
                ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage,
                ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage>(
                  this, METHODID_EXECUTE)))
          .addMethod(
            getInsertMethodHelper(),
            asyncUnaryCall(
              new MethodHandlers<
                ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertRequest,
                ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertResponse>(
                  this, METHODID_INSERT)))
          .addMethod(
            getRemoveMethodHelper(),
            asyncUnaryCall(
              new MethodHandlers<
                ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveRequest,
                ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveResponse>(
                  this, METHODID_REMOVE)))
          .addMethod(
            getQueryByKeyMethodHelper(),
            asyncServerStreamingCall(
              new MethodHandlers<
                ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest,
                ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse>(
                  this, METHODID_QUERY_BY_KEY)))
          .addMethod(
            getQueryByFilterMethodHelper(),
            asyncServerStreamingCall(
              new MethodHandlers<
                ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest,
                ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse>(
                  this, METHODID_QUERY_BY_FILTER)))
          .addMethod(
            getTransferParameterMethodHelper(),
            asyncUnaryCall(
              new MethodHandlers<
                ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaRequest,
                ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaResponse>(
                  this, METHODID_TRANSFER_PARAMETER)))
          .addMethod(
            getPatchNrfProfileMethodHelper(),
            asyncUnaryCall(
              new MethodHandlers<
                ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchRequest,
                ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchResponse>(
                  this, METHODID_PATCH_NRF_PROFILE)))
          .build();
    }
  }

  /**
   */
  public static final class NFDataManagementServiceStub extends io.grpc.stub.AbstractStub<NFDataManagementServiceStub> {
    private NFDataManagementServiceStub(io.grpc.Channel channel) {
      super(channel);
    }

    private NFDataManagementServiceStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected NFDataManagementServiceStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new NFDataManagementServiceStub(channel, callOptions);
    }

    /**
     */
    public void execute(ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage request,
        io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getExecuteMethodHelper(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void insert(ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertRequest request,
        io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getInsertMethodHelper(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void remove(ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveRequest request,
        io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getRemoveMethodHelper(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void queryByKey(ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest request,
        io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse> responseObserver) {
      asyncServerStreamingCall(
          getChannel().newCall(getQueryByKeyMethodHelper(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void queryByFilter(ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest request,
        io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse> responseObserver) {
      asyncServerStreamingCall(
          getChannel().newCall(getQueryByFilterMethodHelper(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void transferParameter(ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaRequest request,
        io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getTransferParameterMethodHelper(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void patchNrfProfile(ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchRequest request,
        io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getPatchNrfProfileMethodHelper(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class NFDataManagementServiceBlockingStub extends io.grpc.stub.AbstractStub<NFDataManagementServiceBlockingStub> {
    private NFDataManagementServiceBlockingStub(io.grpc.Channel channel) {
      super(channel);
    }

    private NFDataManagementServiceBlockingStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected NFDataManagementServiceBlockingStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new NFDataManagementServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage execute(ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage request) {
      return blockingUnaryCall(
          getChannel(), getExecuteMethodHelper(), getCallOptions(), request);
    }

    /**
     */
    public ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertResponse insert(ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertRequest request) {
      return blockingUnaryCall(
          getChannel(), getInsertMethodHelper(), getCallOptions(), request);
    }

    /**
     */
    public ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveResponse remove(ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveRequest request) {
      return blockingUnaryCall(
          getChannel(), getRemoveMethodHelper(), getCallOptions(), request);
    }

    /**
     */
    public java.util.Iterator<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse> queryByKey(
        ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest request) {
      return blockingServerStreamingCall(
          getChannel(), getQueryByKeyMethodHelper(), getCallOptions(), request);
    }

    /**
     */
    public java.util.Iterator<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse> queryByFilter(
        ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest request) {
      return blockingServerStreamingCall(
          getChannel(), getQueryByFilterMethodHelper(), getCallOptions(), request);
    }

    /**
     */
    public ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaResponse transferParameter(ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaRequest request) {
      return blockingUnaryCall(
          getChannel(), getTransferParameterMethodHelper(), getCallOptions(), request);
    }

    /**
     */
    public ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchResponse patchNrfProfile(ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchRequest request) {
      return blockingUnaryCall(
          getChannel(), getPatchNrfProfileMethodHelper(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class NFDataManagementServiceFutureStub extends io.grpc.stub.AbstractStub<NFDataManagementServiceFutureStub> {
    private NFDataManagementServiceFutureStub(io.grpc.Channel channel) {
      super(channel);
    }

    private NFDataManagementServiceFutureStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected NFDataManagementServiceFutureStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new NFDataManagementServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage> execute(
        ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage request) {
      return futureUnaryCall(
          getChannel().newCall(getExecuteMethodHelper(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertResponse> insert(
        ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getInsertMethodHelper(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveResponse> remove(
        ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getRemoveMethodHelper(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaResponse> transferParameter(
        ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getTransferParameterMethodHelper(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchResponse> patchNrfProfile(
        ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getPatchNrfProfileMethodHelper(), getCallOptions()), request);
    }
  }

  private static final int METHODID_EXECUTE = 0;
  private static final int METHODID_INSERT = 1;
  private static final int METHODID_REMOVE = 2;
  private static final int METHODID_QUERY_BY_KEY = 3;
  private static final int METHODID_QUERY_BY_FILTER = 4;
  private static final int METHODID_TRANSFER_PARAMETER = 5;
  private static final int METHODID_PATCH_NRF_PROFILE = 6;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final NFDataManagementServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(NFDataManagementServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_EXECUTE:
          serviceImpl.execute((ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage) request,
              (io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage>) responseObserver);
          break;
        case METHODID_INSERT:
          serviceImpl.insert((ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertRequest) request,
              (io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertResponse>) responseObserver);
          break;
        case METHODID_REMOVE:
          serviceImpl.remove((ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveRequest) request,
              (io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveResponse>) responseObserver);
          break;
        case METHODID_QUERY_BY_KEY:
          serviceImpl.queryByKey((ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest) request,
              (io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse>) responseObserver);
          break;
        case METHODID_QUERY_BY_FILTER:
          serviceImpl.queryByFilter((ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest) request,
              (io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse>) responseObserver);
          break;
        case METHODID_TRANSFER_PARAMETER:
          serviceImpl.transferParameter((ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaRequest) request,
              (io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaResponse>) responseObserver);
          break;
        case METHODID_PATCH_NRF_PROFILE:
          serviceImpl.patchNrfProfile((ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchRequest) request,
              (io.grpc.stub.StreamObserver<ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchResponse>) responseObserver);
          break;
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class NFDataManagementServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    NFDataManagementServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("NFDataManagementService");
    }
  }

  private static final class NFDataManagementServiceFileDescriptorSupplier
      extends NFDataManagementServiceBaseDescriptorSupplier {
    NFDataManagementServiceFileDescriptorSupplier() {}
  }

  private static final class NFDataManagementServiceMethodDescriptorSupplier
      extends NFDataManagementServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    NFDataManagementServiceMethodDescriptorSupplier(String methodName) {
      this.methodName = methodName;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.MethodDescriptor getMethodDescriptor() {
      return getServiceDescriptor().findMethodByName(methodName);
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (NFDataManagementServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new NFDataManagementServiceFileDescriptorSupplier())
              .addMethod(getExecuteMethodHelper())
              .addMethod(getInsertMethodHelper())
              .addMethod(getRemoveMethodHelper())
              .addMethod(getQueryByKeyMethodHelper())
              .addMethod(getQueryByFilterMethodHelper())
              .addMethod(getTransferParameterMethodHelper())
              .addMethod(getPatchNrfProfileMethodHelper())
              .build();
        }
      }
    }
    return result;
  }
}

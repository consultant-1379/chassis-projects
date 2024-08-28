/*package ericsson.core.nrf.dbproxy.dbproxy_client;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.StatusRuntimeException;
import java.util.concurrent.TimeUnit;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceGrpc;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;

public class DBProxyClient
{
    private final NFDataManagementServiceGrpc.NFDataManagementServiceBlockingStub blockingStub;
    private final ManagedChannel channel;

    public DBProxyClient(String host, int port)
    {
	channel = ManagedChannelBuilder.forAddress(host, port).usePlaintext().build();
	blockingStub = NFDataManagementServiceGrpc.newBlockingStub(channel);
    }

    public void send()
    {
	try
	{
	    NFMessage request = RequestBuilder.buildNFProfileGetRequest();
	    System.out.println(request.toString());
	    NFMessage response = blockingStub.execute(request);
	    System.out.println(response.toString());
	}
	catch(StatusRuntimeException e)
	{
	    System.out.println(e.toString());
	}
    }

    public void shutdown() throws InterruptedException {
    	channel.shutdown().awaitTermination(5, TimeUnit.SECONDS);
    }

    public static void main(String[] args) throws Exception
    {
	System.out.println("This is DBProxyClient");
	DBProxyClient client = new DBProxyClient("localhost", 50051);

	try
	{
	    client.send();
	}
	finally
	{
	   client.shutdown();
	}
    }
}
*/

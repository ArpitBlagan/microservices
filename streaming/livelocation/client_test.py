import grpc
import location_pb2
import location_pb2_grpc

def run():
    channel = grpc.insecure_channel('localhost:50051')
    stub = location_pb2_grpc.LocationServiceStub(channel)
    def location_requests():
        locations = [
            location_pb2.LocationRequest(user_id=1, latitude=37.7749, longitude=-122.4194),
            location_pb2.LocationRequest(user_id=2, latitude=40.7128, longitude=-74.0060),
            location_pb2.LocationRequest(user_id=3, latitude=34.0522, longitude=-118.2437),
        ]
        for location in locations:
            yield location
    try:
        response = stub.StreamLocation(location_requests())
        for res in response:
            print(f"res :{res}")
    except grpc.RpcError as e:
        print(f"RPC error: {e.details()}")
if __name__ == "__main__":
    run()
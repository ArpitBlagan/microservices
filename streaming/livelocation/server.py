import grpc
from concurrent import futures
import location_pb2
import location_pb2_grpc
class LocationService(location_pb2_grpc.LocationServiceServicer):
    def StreamLocation(self, request_iterator, context):
        try:
            # Iterate over the incoming stream of requests
            for location in request_iterator:
                print(f"Received location: {location}")
                # Process the location and send a response
                yield location_pb2.LocationResponse(status="received")
        except Exception as e:
            # Log the exception and send an error response
            context.set_details(f"Error processing request: {str(e)}")
            context.set_code(grpc.StatusCode.INTERNAL)
            raise e
def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    location_pb2_grpc.add_LocationServiceServicer_to_server(LocationService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("gRPC Server started on port 50051")
    server.wait_for_termination()

if __name__ == "__main__":
    serve()
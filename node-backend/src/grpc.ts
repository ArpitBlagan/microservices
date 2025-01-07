import grpc from "@grpc/grpc-js";
import protoLoader from "@grpc/proto-loader";
import { sendEmail } from "./controllers/grpc";

const PROTO_PATH = "./protoc/crud.proto";
const packageDefinition = protoLoader.loadSync(PROTO_PATH);
const grpcObject = grpc.loadPackageDefinition(packageDefinition);
//@ts-ignore
const emailService = grpcObject.Email.CrudService;

const grpcServer = new grpc.Server();
grpcServer.addService(emailService.service, { SendEmail: sendEmail });

grpcServer.bindAsync(
  "127.0.0.1:50051",
  grpc.ServerCredentials.createInsecure(),
  () => {
    console.log("gRPC server running on port 50051");
    grpcServer.start();
  }
);

import express from "express";
import cors from "cors";
import cookieParser from "cookie-parser";
import { router } from "./router";
import dotenv from "dotenv";
dotenv.config();
//This Backend will handle sending notification part for the app throught email, sms etc.
const app = express();

app.use(
  cors({
    origin: [],
    credentials: true,
  })
);
app.use(express.json());
app.use(cookieParser());

app.use("/api", router);

app.listen(process.env.PORT, () => {
  console.log(`Server listening on port ${process.env.PORT}`);
});

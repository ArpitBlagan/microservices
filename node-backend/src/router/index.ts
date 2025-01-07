import { Router } from "express";
import { sendEmail } from "../controllers";

export const router = Router();

router.route("/sendEmail").post(sendEmail);

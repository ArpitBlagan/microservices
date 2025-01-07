import { Request, Response } from "express";
import { rsendEmail } from "../resend";

export async function sendEmail(req: Request, res: Response) {
  const body = req.body;
  try {
    // use resend to send the email...
    await rsendEmail({});
    res.status(200).json({ message: "Email send sucessfully :)" });
  } catch (err) {
    res.status(500).json({ message: "Internal server error :(" });
  }
}

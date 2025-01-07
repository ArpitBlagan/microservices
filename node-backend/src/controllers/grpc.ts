import { rsendEmail } from "../resend";

export async function sendEmail(call: any, callback: any) {
  const body = call.req;
  try {
    // use resend to send the email...
    await rsendEmail({});
    callback({ message: "Email send sucessfully :)" });
  } catch (err) {
    callback({ message: "Internal server error :(" });
  }
}

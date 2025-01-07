import { Resend } from "resend";
import dotenv from "dotenv";
dotenv.config();
const resend = new Resend(process.env.API_KEY as string);

export const rsendEmail = async (userInfo: any) => {
  const { data, error } = await resend.emails.send({
    from: `Arpit Blagan <onboarding@resend.dev>`,
    to: [`blaganarpit@gmail.com`],
    subject: `Welcome to our App.`,
    html: `<p>Thank you for registering on our app.</p>`,
  });

  if (error) {
    console.log(error);
    throw Error("Not able to send email may be its Internal server error :(");
  }
  console.log(data);
};

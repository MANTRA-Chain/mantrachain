import * as dotenv from 'dotenv'
import { setup } from "./src/helpers/env"

dotenv.config()

export default async () => {
  const host = process.env.API_URL || 'http://localhost:1317';
  await setup(host);
};

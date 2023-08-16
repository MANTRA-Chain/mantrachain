import { setup } from "./src/helpers/env"

export default async () => {
  const host = process.env.API_URL || 'http://localhost:1317';
  await setup(host);
};

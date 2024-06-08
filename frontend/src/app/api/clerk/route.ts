import { fetchBackend } from '@/config/db';
import { WebhookEvent } from '@clerk/nextjs/server'
import { headers } from 'next/headers'
import { NextRequest, NextResponse } from 'next/server'
import { Webhook } from 'svix'

interface CreateUserInput {
  email: string;
  name: string;
  display_name: string;
  balance?: number;
}

interface UpdateUserInput {
  email?: string;
  name?: string;
  display_name?: string;
  balance?: number;
}

const webhookSecret = process.env.CLERK_WEBHOOK_SECRET || ``

async function validateRequest(request: NextRequest) {
  const payloadString = await request.text()
  const headerPayload = headers()

  const svixHeaders = {
    'svix-id': headerPayload.get('svix-id')!,
    'svix-timestamp': headerPayload.get('svix-timestamp')!,
    'svix-signature': headerPayload.get('svix-signature')!,
  }
  const wh = new Webhook(webhookSecret)
  return wh.verify(payloadString, svixHeaders) as WebhookEvent
}

export async function POST(request: NextRequest) {
  try {
    const payload = await validateRequest(request);
    const data: any = payload.data;
    console.log(payload)
    if (payload.type === 'user.created') {
      await createUser({
        email: data["email_addresses"][0]["email_address"],
        name: data["first_name"],
        display_name: data["username"] || data["email_addresses"][0]["email_address"],
        balance: 0.01
      })
    } else if (payload.type === 'user.updated') {
      await updateUser(data["email_addresses"][0]["email_address"], {
        display_name: data["username"]
      })
    }

    return NextResponse.json({ message: 'User updated' })
  } catch (error) {
    console.error(error)
    return NextResponse.json({ error: 'Invalid request' }, { status: 400 })
  }
}

async function createUser(input: CreateUserInput) {
  try {
    const response = await fetchBackend<{ message: string }>(
      {endpoint: '/users', method:  'POST', data: input}
    );
    return response.data;
  } catch (error) {
    console.error(error);
  }
}

async function updateUser(email: string, input: UpdateUserInput) {
  try {
    const response = await fetchBackend<{ message: string }>(
      {endpoint: `/users/${email}`, method:  'PUT', data: input}
    );
    return response.data;
  } catch(error) {
    console.error("Error updating user");
  }
}
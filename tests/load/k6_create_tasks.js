import http from 'k6/http';
import { sleep } from 'k6';

// Define the base URL of your API
const BASE_URL = 'http://localhost:3333';

// Define the options for the load test
export let options = {
  stages: [
    { duration: '1m', target: 50 },  // Ramp-up to 50 users over 1 minute
    { duration: '1m', target: 50 },  // Stay at 50 users for 3 minutes
    { duration: '1m', target: 0 },   // Ramp-down to 0 users over 1 minute
  ],
};

// Define the login function
export function login() {
  const url = `${BASE_URL}/login`;
  const payload = JSON.stringify({
    username: 'admin',
    password: 'some_password',
  });
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const response = http.post(url, payload, params);
  return response.json();
}

// Define the task creation function
export function createTask(authToken) {
  const url = `${BASE_URL}/tasks`;
  const payload = JSON.stringify({
    userId: '1',
    description: 'Task 1'
  });
  const params = {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${authToken}`,
    },
  };

  const response = http.post(url, payload, params);
  
  // Log the response status and body
  console.log(`Create Task Response: ${response.status}, ${response.body}`);

  // Optionally, check the response status and handle errors
  if (response.status !== 201) {
    console.error(`Failed to create task. Status code: ${response.status}`);
  }
}

// Main function to execute the load test
export default function () {
  const authToken = login().token; // Assuming your login response returns a token
  createTask(authToken);
  sleep(1);
}

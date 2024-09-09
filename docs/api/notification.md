# Notification Service Documentation

## Notification API

The notification service handles both email and in-app notifications. It uses RabbitMQ to ensure reliable delivery of notifications, even during high load or temporary service disruptions.

### SendNotification

Endpoint: `NotificationService.SendNotification`
Authorization: Authenticated users

Request:
```json
{
  "type": "EMAIL",
  "to": "user@example.com",
  "subject": "Welcome to iRankHub",
  "content": "Thank you for joining iRankHub!",
  "token": "your_auth_token_here"
}
```

### GetUnreadNotifications

Endpoint: `NotificationService.GetUnreadNotifications`
Authorization: Authenticated users

Request:
```json
{
  "user_id": 1,
  "token": "your_auth_token_here"
}
```

### MarkNotificationsAsRead

Endpoint: `NotificationService.MarkNotificationsAsRead`
Authorization: Authenticated users

Request:
```json
{
  "user_id": 1,
  "token": "your_auth_token_here"
}
```

### SubscribeToNotification

Endpoint: `NotificationService.SubscribeToNotification`
Authorization: Authenticated users

Request:
```json
{
  "user_id": 1,
}
```

## Frontend Integration

To integrate notifications in the frontend using gRPC:

1. Set up a gRPC-Web client for the NotificationService.

2. Implement a function to send notifications:

```javascript
async function sendNotification(type, to, subject, content) {
  try {
    const request = new SendNotificationRequest();
    request.setType(type);
    request.setTo(to);
    request.setSubject(subject);
    request.setContent(content);

    const response = await notificationClient.sendNotification(request);
    console.log('Notification sent successfully:', response.getSuccess());
  } catch (error) {
    console.error('Failed to send notification:', error);
  }
}
```

3. Implement a notification center component that fetches and displays unread notifications:

```jsx
function NotificationCenter() {
  const [notifications, setNotifications] = useState([]);

  useEffect(() => {
    fetchUnreadNotifications();
  }, []);

  async function fetchUnreadNotifications() {
    const request = new GetUnreadNotificationsRequest();
    request.setUserId(currentUserId);

    try {
      const response = await notificationClient.getUnreadNotifications(request);
      setNotifications(response.getNotificationsList());
    } catch (error) {
      console.error('Failed to fetch unread notifications:', error);
    }
  }

  return (
    <div className="notification-center">
      {notifications.map((notification) => (
        <NotificationItem key={notification.getId()} notification={notification} />
      ))}
    </div>
  );
}
```

4. Implement real-time notifications using gRPC streaming:

```javascript
function subscribeToNotifications(userId) {
  const request = new SubscribeRequest();
  request.setUserId(userId);

  const stream = notificationClient.subscribeToNotifications(request);

  stream.on('data', (response) => {
    const newNotification = response.getNotification();
    // Handle the new notification (e.g., update UI, show a toast message)
    console.log('New notification received:', newNotification);
  });

  stream.on('error', (error) => {
    console.error('Error in notification stream:', error);
    // Implement reconnection logic here
  });

  stream.on('end', () => {
    console.log('Notification stream ended');
    // Implement reconnection logic here
  });
}
```

Call `subscribeToNotifications(currentUserId)` when the user logs in or the application starts.

## Testing Notification Features

To test the notification features:

1. Ensure your gRPC client is properly set up and connected to the backend.

2. Test the following scenarios:

   a. Sending Notifications:
   - Use `SendNotification` to send an email notification.
   - Use `SendNotification` to send an in-app notification.
   - Verify that the notifications are received (check email and frontend notification center).

   b. Retrieving Notifications:
   - Use `GetUnreadNotifications` to fetch unread notifications for a user.
   - Verify that the frontend displays the unread notifications correctly.

   c. Marking Notifications as Read:
   - Use `MarkNotificationsAsRead` to mark all notifications as read for a user.
   - Use `GetUnreadNotifications` again to verify that there are no unread notifications.

   d. Real-time Notifications:
   - Implement the `SubscribeToNotifications` streaming RPC in your frontend.
   - Open multiple browser windows or tabs with the frontend application.
   - Trigger an action that sends a notification (e.g., tournament invitation).
   - Verify that the notification appears in real-time across all open sessions for the recipient.

3. For each test, verify that:
   - Email notifications are sent and received correctly.
   - In-app notifications appear in the frontend notification center.
   - Real-time updates work as expected using the gRPC stream.

4. Test error handling:
   - Try sending a notification with invalid parameters.
   - Verify that appropriate error messages are displayed.

5. Performance testing:
   - Send a large number of notifications in quick succession.
   - Verify that all notifications are processed and delivered correctly.

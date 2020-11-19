from notification_providers.telegram_notification import TelegramNotifier

if __name__ == '__main__':
    tn = TelegramNotifier()
    tn.start_polling()

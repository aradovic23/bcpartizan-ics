# 🚀 Deployment Checklist

## Quick Steps for You

### 1️⃣ Create GitHub Repository (2 minutes)

1. Go to https://github.com/new
2. Repository name: `bcpartizan-ics` (or any name)
3. Choose **Public** or **Private** (your choice)
4. **DON'T** click "Add a README" (we have files already)
5. Click **Create repository**
6. **Copy the repository URL** from the page (looks like: `https://github.com/YOUR_USERNAME/bcpartizan-ics.git`)

### 2️⃣ Push Code to GitHub (1 minute)

Open Terminal and run these commands:

```bash
cd /Users/user_name/dev/bcpartizan-ics

# Add your GitHub URL (REPLACE with your actual URL from step 1)
git remote add origin https://github.com/YOUR_USERNAME/bcpartizan-ics.git

# Stage all files
git add .

# Commit
git commit -m "Initial commit: Partizan Basketball ICS Calendar"

# Push to GitHub
git branch -M main
git push -u origin main
```

**Note:** GitHub might ask for authentication. Use a personal access token if needed.

### 4️⃣ Get Your Calendar URL

Once deployed:
- Your service URL will look like: `https://partizan-ics-calendar.onrender.com`
- Your calendar subscription URL is: `https://partizan-ics-calendar.onrender.com/calendar.ics`

**Share this URL with friends!**

---

## 📋 What to Share with Friends

Copy and send this to your friends:

```
🏀 KK Partizan Basketball Calendar

Subscribe to get all Partizan games in your calendar!

📅 Calendar URL:
https://YOUR-APP-NAME.onrender.com/calendar.ics

How to add:

📱 iPhone/Mac:
• Open Calendar app
• File → New Calendar Subscription
• Paste the URL above

💻 Google Calendar:
• Go to calendar.google.com
• Click "+" next to "Other calendars"
• Select "From URL"
• Paste the URL above

📧 Outlook:
• Add Calendar → Subscribe from web
• Paste the URL above

What you'll get:
✅ All 36 Euroleague games (full season through April 2026)
✅ ABA League games (rolling schedule)
✅ 30-minute reminders before each game
✅ Full venue information
✅ Automatic updates every 2 days

Idemo Partizan! 🖤🤍
```

---

## 🆘 Troubleshooting

### "Service is spinning down"
This is normal on free tier. Service sleeps after 15 min of inactivity.
- First request takes 30-60 seconds
- Calendar apps handle this automatically
- **Optional:** Use UptimeRobot to keep it awake

### "Build failed"
Check Render build logs for errors. Common fixes:
- Make sure all files are pushed to GitHub
- Check `package.json` has correct dependencies
- Verify `render.yaml` is in root directory

### "Can't access calendar"
- Make sure URL ends with `/calendar.ics`
- Check service is "Live" in Render dashboard
- Try accessing in browser first to verify

### Authentication Issues with GitHub
If git push asks for password:
1. Go to https://github.com/settings/tokens
2. Generate new token (classic)
3. Copy token and use as password

---

## ✅ Done!

Everything is ready. Just follow the 4 steps above and you'll have your calendar live in ~6 minutes!

Questions? Check the main README.md or Render documentation.

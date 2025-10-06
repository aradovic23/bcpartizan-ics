# ⚡ Quick Start - Deploy in 5 Minutes

## What I Did For You ✅

✅ Initialized git repository
✅ Created `render.yaml` configuration
✅ Updated `.env.example`
✅ Updated README with deployment instructions
✅ Created `DEPLOYMENT.md` with detailed steps
✅ Everything is ready to deploy!

---

## What YOU Need To Do (5 minutes)

### Step 1: Create GitHub Repository (2 min)

1. Go to https://github.com/new
2. Name: `bcpartizan-ics`
3. Keep public or private (your choice)
4. **DON'T add README** (we have files)
5. Click **Create**
6. **COPY the repository URL** shown on the page

### Step 2: Push to GitHub (1 min)

Open Terminal and run:

```bash
cd /Users/milicacurcic/dev/bcpartizan-ics

# Replace YOUR_USERNAME with your actual GitHub username
git remote add origin https://github.com/YOUR_USERNAME/bcpartizan-ics.git

git add .
git commit -m "Initial commit: Partizan ICS Calendar"
git branch -M main
git push -u origin main
```

### Step 3: Deploy on Render (2 min)

1. Go to https://render.com
2. Click **"Get Started"** (sign in with GitHub)
3. Click **"New +"** → **"Web Service"**
4. Select your `bcpartizan-ics` repository
5. Click **"Create Web Service"**
6. Wait 2-3 minutes ⏱️
7. Done! 🎉

### Step 4: Get Your URL

Your calendar URL will be:
```
https://partizan-ics-calendar.onrender.com/calendar.ics
```

Share this with friends!

---

## Need Help?

📖 **Detailed Instructions:** See `DEPLOYMENT.md`
📚 **Full Documentation:** See `README.md`

---

## Quick Test Before Pushing

Want to test locally first?

```bash
npm start
```

Then open: http://localhost:3000

You should see 49 games (36 Euroleague + 13 ABA League)

---

**Idemo Partizan!** 🖤🤍

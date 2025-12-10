@echo off
echo ========================================
echo  Starting RestySched with n8n
echo ========================================
echo.
echo This will start:
echo  - MongoDB (port 27017)
echo  - n8n (port 5678)
echo  - RestySched (port 8080)
echo.
echo After starting, open:
echo  - RestySched: http://localhost:8080
echo  - n8n:        http://localhost:5678 (admin/admin)
echo.
echo Press any key to start...
pause >nul

docker-compose -f docker-compose.n8n.yml up -d

echo.
echo ========================================
echo  Services started!
echo ========================================
echo.
echo Opening RestySched in browser...
timeout /t 3 >nul
start http://localhost:8080

echo.
echo Opening n8n in browser...
timeout /t 2 >nul
start http://localhost:5678

echo.
echo ========================================
echo  Next Steps:
echo ========================================
echo  1. Set up webhook in n8n (see N8N_SETUP.md)
echo  2. Copy webhook URL
echo  3. Add to .env: N8N_WEBHOOK_URL=your-webhook-url
echo  4. Restart: docker-compose -f docker-compose.n8n.yml restart app
echo.
echo To stop all services:
echo   docker-compose -f docker-compose.n8n.yml down
echo.
pause

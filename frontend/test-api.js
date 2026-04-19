// Simple API test script to verify frontend-backend integration
const axios = require('axios');

const API_BASE_URL = 'http://localhost:8080/api';

async function testAPI() {
  console.log('Testing API integration...');
  
  try {
    // Test login
    console.log('1. Testing login...');
    const loginResponse = await axios.post(`${API_BASE_URL}/auth/login`, {
      email: 'iiceekiingfx@gmail.com',
      password: 'Zephyr@651818'
    });
    
    const { access_token } = loginResponse.data;
    console.log('   Login successful! Token received.');
    
    // Set up auth headers
    const authHeaders = {
      Authorization: `Bearer ${access_token}`
    };
    
    // Test courses endpoint
    console.log('2. Testing courses endpoint...');
    const coursesResponse = await axios.get(`${API_BASE_URL}/courses/`, { headers: authHeaders });
    console.log(`   Found ${coursesResponse.data.count} courses.`);
    
    // Test signals endpoint
    console.log('3. Testing signals endpoint...');
    const signalsResponse = await axios.get(`${API_BASE_URL}/signals/`, { headers: authHeaders });
    console.log(`   Found ${signalsResponse.data.signals.length} signals.`);
    
    // Test portfolio endpoint
    console.log('4. Testing portfolio endpoint...');
    const portfolioResponse = await axios.get(`${API_BASE_URL}/portfolio/accounts`, { headers: authHeaders });
    console.log(`   Found ${portfolioResponse.data.accounts.length} trading accounts.`);
    
    // Test dashboard endpoint
    console.log('5. Testing dashboard endpoint...');
    const dashboardResponse = await axios.get(`${API_BASE_URL}/dashboard/overview`, { headers: authHeaders });
    console.log(`   Dashboard data: ${JSON.stringify(dashboardResponse.data, null, 2)}`);
    
    console.log('\nAll API tests passed! Frontend-backend integration is working correctly.');
    
  } catch (error) {
    console.error('API test failed:', error.response?.data || error.message);
    process.exit(1);
  }
}

testAPI();

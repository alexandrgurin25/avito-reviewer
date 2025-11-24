import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend, Rate, Counter } from 'k6/metrics';

// Кастомные метрики
const errorCount = new Counter('errors');
const prCreationTime = new Trend('pr_creation_time');
const reassignmentTime = new Trend('reassignment_time');
const successRate = new Rate('success_rate');

export const options = {
  stages: [
   { duration: '2m', target: 170 },   // 50 concurrent users
    { duration: '3m', target: 200 },  // 200 concurrent users  
    { duration: '2m', target: 120 },  // 150 concurrent users
  ],
  thresholds: {
    http_req_duration: ['p(95)<300'],
    errors: ['count<100'],
    success_rate: ['rate>0.999'],
  },
};

const BASE_URL = 'http://localhost:8080';

// Тестовые данные - 10 команд
const TEAMS = [
  'backend-loadtest', 'frontend-loadtest', 'payments-loadtest', 'mobile-loadtest', 
  'data-science-loadtest', 'devops-loadtest', 'qa-loadtest', 
  'security-loadtest', 'infrastructure-loadtest', 'analytics-loadtest'
];

// Функция для создания пользователей команды
function createTeamUsers(team) {
  const users = [];
  for (let i = 1; i <= 10; i++) {
    users.push({
      user_id: `user-${team}-${i}`,
      username: `User ${i} (${team})`,
      is_active: true,
    });
  }
  return users;
}

export function setup() {
  console.log('Setting up test data for 10 teams...');
  
  const createdTeams = [];
  
  // Создаем команды с уникальными именами для теста
  TEAMS.forEach(team => {
    const users = createTeamUsers(team);
    const payload = JSON.stringify({
      team_name: team,
      members: users,
    });

    const res = http.post(`${BASE_URL}/team/add`, payload, {
      headers: { 'Content-Type': 'application/json' },
    });

    if (res.status === 201) {
      console.log(`✅ Team ${team} created successfully`);
      createdTeams.push({
        name: team,
        users: users
      });
    } else if (res.status === 400) {
      console.log(`⚠️ Team ${team} already exists, using existing team`);
      // Если команда уже существует, всё равно добавляем её для тестирования
      createdTeams.push({
        name: team,
        users: users
      });
    } else {
      console.log(`❌ Failed to create team ${team}: ${res.status}`);
    }
  });
  
  console.log(`Setup complete. ${createdTeams.length} teams ready for testing`);
  return { teams: createdTeams };
}

export default function(data) {
  if (!data.teams || data.teams.length === 0) {
    console.log('No teams available for testing');
    return;
  }

  const randomTeam = data.teams[Math.floor(Math.random() * data.teams.length)];
  const teamUsers = randomTeam.users.filter(u => u.is_active);
  
  if (teamUsers.length === 0) {
    errorCount.add(1);
    return;
  }

  const user = teamUsers[Math.floor(Math.random() * teamUsers.length)];
  const prId = `pr-${Date.now()}-${Math.random().toString(36).substr(2, 5)}`;

  // Определяем тип операции для разнообразия тестов
  const operationType = Math.random();
  
  if (operationType < 0.6) {
    // 60%: Создание нового PR
    testCreatePR(randomTeam.name, user, prId);
  } else if (operationType < 0.8) {
    // 20%: Получение информации о команде
    testGetTeam(randomTeam.name);
  } else {
    // 20%: Получение PR пользователя
    testGetUserPRs(user);
  }

  sleep(0.5);
}

function testCreatePR(team, user, prId) {
  const payload = JSON.stringify({
    pull_request_id: prId,
    pull_request_name: `Feature for ${team} - ${prId}`,
    author_id: user.user_id,
  });

  const startTime = Date.now();
  const res = http.post(`${BASE_URL}/pullRequest/create`, payload, {
    headers: { 'Content-Type': 'application/json' },
  });
  const duration = Date.now() - startTime;

  const success = check(res, {
    'PR created successfully': (r) => r.status === 201,
    'PR has reasonable response time': (r) => duration < 1000,
  });

  if (res.status === 201) {
    prCreationTime.add(duration);
  }

  successRate.add(success);
  if (!success && res.status !== 409) { // 409 - PR уже существует, это нормально
    errorCount.add(1);
    console.log(`❌ PR creation failed: ${res.status} - ${res.body}`);
  }
}

function testGetTeam(team) {
  const res = http.get(`${BASE_URL}/team/get?team_name=${team}`);
  
  check(res, {
    'team retrieved successfully': (r) => r.status === 200,
  });
}

function testGetUserPRs(user) {
  const res = http.get(`${BASE_URL}/users/getReview?user_id=${user.user_id}`);
  
  check(res, {
    'user PRs retrieved successfully': (r) => r.status === 200,
  });
}

export function teardown() {
  console.log('🎉 Load test completed!');
}
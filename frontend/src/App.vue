<script setup>
import { ref, reactive } from 'vue'
import axios from 'axios'

// ตัวเลือกสำหรับ Dropdown Occupation
const occupations = [
  'occupations 1',
  'occupations 2',
  'occupations 3'
]

// ข้อมูลฟอร์ม
const form = reactive({
  first_name: '',
  last_name: '',
  email: '',
  phone: '',
  profile_base64: '',
  birth_day: '', // เก็บค่า YYYY-MM-DD จากปฏิทิน
  occupation: '',
  sex: 'Male'
})

// เก็บข้อความ Error ของแต่ละฟิลด์
const errors = reactive({
  first_name: '',
  last_name: '',
  email: '',
  phone: '',
  profile_base64: '',
  birth_day: '',
  occupation: '',
  sex: ''
})

const fileInput = ref(null)
const fileName = ref('')
const statusMessage = ref('')
const isErrorStatus = ref(false)
const isSubmitting = ref(false)

// แปลงไฟล์รูปภาพที่อัปโหลดเป็น Base64
const handleFileUpload = (event) => {
  const file = event.target.files[0]
  if (!file) {
    form.profile_base64 = ''
    fileName.value = ''
    return
  }

  fileName.value = file.name
  const reader = new FileReader()
  reader.onload = (e) => {
    form.profile_base64 = e.target.result
    errors.profile_base64 = ''
  }
  reader.readAsDataURL(file)
}

// ตรวจสอบความถูกต้องของข้อมูลฝั่ง Frontend
const validateForm = () => {
  let isValid = true
  Object.keys(errors).forEach((key) => (errors[key] = ''))

  if (!form.first_name.trim()) {
    errors.first_name = 'Please enter First Name'
    isValid = false
  }
  if (!form.last_name.trim()) {
    errors.last_name = 'Please enter Last Name'
    isValid = false
  }

  // ตรวจ Format Email
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!form.email.trim()) {
    errors.email = 'Please enter Email'
    isValid = false
  } else if (!emailRegex.test(form.email)) {
    errors.email = 'Invalid Email format'
    isValid = false
  }

  // ตรวจ Format เบอร์โทร (9-10 หลัก ขึ้นต้นด้วย 0)
  const phoneRegex = /^0[0-9]{8,9}$/
  if (!form.phone.trim()) {
    errors.phone = 'Please enter Phone'
    isValid = false
  } else if (!phoneRegex.test(form.phone)) {
    errors.phone = 'Phone must be 9-10 digits starting with 0'
    isValid = false
  }

  if (!form.profile_base64) {
    errors.profile_base64 = 'Please upload Profile image'
    isValid = false
  }

  if (!form.birth_day) {
    errors.birth_day = 'Please select Birth Day'
    isValid = false
  }

  if (!form.occupation) {
    errors.occupation = 'Please select Occupation'
    isValid = false
  }

  if (!form.sex) {
    errors.sex = 'Please select Sex'
    isValid = false
  }

  return isValid
}

// กดปุ่ม Save เพื่อส่งข้อมูลไปหา Go Backend
const handleSave = async () => {
  statusMessage.value = ''
  if (!validateForm()) return

  // แปลงรูปแบบวันที่จาก YYYY-MM-DD เป็น DD/MM/YYYY ก่อนส่งให้ Backend ตามที่ตารางต้องการ
  let formattedDate = form.birth_day
  if (form.birth_day) {
    const parts = form.birth_day.split('-') // [YYYY, MM, DD]
    if (parts.length === 3) {
      formattedDate = `${parts[2]}/${parts[1]}/${parts[0]}`
    }
  }

  const payload = {
    ...form,
    birth_day: formattedDate
  }

  isSubmitting.value = true
  try {
    const response = await axios.post('http://localhost:8080/api/profiles', payload)
    isErrorStatus.value = false
    statusMessage.value = `${response.data.message} (ID: ${response.data.id})`
    handleClear(false)
  } catch (err) {
    isErrorStatus.value = true
    statusMessage.value = err.response?.data?.error || 'Failed to connect to server'
  } finally {
    isSubmitting.value = false
  }
}

// กดปุ่ม Clear เพื่อล้างค่าทั้งหมดในฟอร์ม
const handleClear = (clearStatus = true) => {
  form.first_name = ''
  form.last_name = ''
  form.email = ''
  form.phone = ''
  form.profile_base64 = ''
  form.birth_day = ''
  form.occupation = ''
  form.sex = 'Male'

  fileName.value = ''
  if (fileInput.value) {
    fileInput.value.value = ''
  }

  Object.keys(errors).forEach((key) => (errors[key] = ''))
  if (clearStatus) {
    statusMessage.value = ''
  }
}
</script>

<template>
  <div class="page-wrapper">
    <div class="form-container">
      <div class="header-bar">IT 04-1</div>

      <form @submit.prevent="handleSave" class="form-grid">
        <!-- แถว 1: First Name / Last Name -->
        <div class="form-group">
          <label>First Name</label>
          <input v-model="form.first_name" type="text" />
          <span v-if="errors.first_name" class="error-text">{{ errors.first_name }}</span>
        </div>

        <div class="form-group">
          <label>Last Name</label>
          <input v-model="form.last_name" type="text" />
          <span v-if="errors.last_name" class="error-text">{{ errors.last_name }}</span>
        </div>

        <!-- แถว 2: Email / Phone -->
        <div class="form-group">
          <label>Email</label>
          <input v-model="form.email" type="text" placeholder="example@mail.com" />
          <span v-if="errors.email" class="error-text">{{ errors.email }}</span>
        </div>

        <div class="form-group">
          <label>Phone</label>
          <input v-model="form.phone" type="text" maxlength="10" placeholder="0812345678" />
          <span v-if="errors.phone" class="error-text">{{ errors.phone }}</span>
        </div>

        <!-- แถว 3: Profile / Birth Day (เปลี่ยนเป็นปฏิทิน) -->
        <div class="form-group">
          <label>Profile</label>
          <div class="file-upload-box">
            <input
              ref="fileInput"
              type="file"
              accept="image/*"
              @change="handleFileUpload"
              id="profile-upload"
              hidden
            />
            <span class="file-name">{{ fileName || 'No file selected' }}</span>
            <label for="profile-upload" class="browse-btn">Browse</label>
          </div>
          <span v-if="errors.profile_base64" class="error-text">{{ errors.profile_base64 }}</span>
        </div>

        <div class="form-group">
          <label>Birth Day</label>
          <input v-model="form.birth_day" type="date" />
          <span v-if="errors.birth_day" class="error-text">{{ errors.birth_day }}</span>
        </div>

        <!-- แถว 4: Occupation -->
        <div class="form-group">
          <label>Occupation</label>
          <select v-model="form.occupation">
            <option value="" disabled>-- Select Occupation --</option>
            <option v-for="occ in occupations" :key="occ" :value="occ">
              {{ occ }}
            </option>
          </select>
          <span v-if="errors.occupation" class="error-text">{{ errors.occupation }}</span>
        </div>

        <div class="form-group empty-col"></div>

        <!-- แถว 5: Sex -->
        <div class="form-group full-width">
          <label>Sex</label>
          <div class="radio-group">
            <label class="radio-label">
              <input type="radio" value="Male" v-model="form.sex" /> Male
            </label>
            <label class="radio-label">
              <input type="radio" value="Female" v-model="form.sex" /> Female
            </label>
          </div>
          <span v-if="errors.sex" class="error-text">{{ errors.sex }}</span>
        </div>

        <!-- ปุ่ม Save และ Clear -->
        <div class="button-row full-width">
          <button type="submit" class="btn btn-save" :disabled="isSubmitting">
            {{ isSubmitting ? 'Saving...' : 'Save' }}
          </button>
          <button type="button" class="btn btn-clear" @click="handleClear(true)">
            Clear
          </button>
        </div>

        <!-- แสดงผลลัพธ์จาก API -->
        <div v-if="statusMessage" :class="['status-banner full-width', isErrorStatus ? 'status-error' : 'status-success']">
          {{ statusMessage }}
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.page-wrapper {
  display: flex;
  justify-content: center;
  padding: 30px 15px;
  font-family: Arial, sans-serif;
  color: #333;
}

.form-container {
  width: 100%;
  max-width: 720px;
  border: 1px solid #ccc;
  background-color: #fff;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
}

.header-bar {
  background-color: #00a65a;
  color: white;
  text-align: center;
  padding: 10px;
  font-weight: bold;
  font-size: 16px;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 18px 24px;
  padding: 25px 30px;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.full-width {
  grid-column: span 2;
}

label {
  font-size: 14px;
  margin-bottom: 6px;
  color: #222;
}

input[type="text"],
input[type="date"],
select {
  height: 36px;
  padding: 0 10px;
  border: 1px solid #888;
  border-radius: 4px;
  font-size: 14px;
  outline: none;
}

input[type="text"]:focus,
input[type="date"]:focus,
select:focus {
  border-color: #00a65a;
}

.file-upload-box {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 36px;
  border: 1px solid #888;
  border-radius: 4px;
  padding-left: 10px;
  overflow: hidden;
}

.file-name {
  font-size: 13px;
  color: #666;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 180px;
}

.browse-btn {
  background-color: #b5b5b5;
  color: #000;
  height: 100%;
  display: flex;
  align-items: center;
  padding: 0 14px;
  margin: 0;
  cursor: pointer;
  font-size: 13px;
  border-left: 1px solid #888;
}

.browse-btn:hover {
  background-color: #a0a0a0;
}

.radio-group {
  display: flex;
  gap: 24px;
  align-items: center;
  margin-top: 4px;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  margin: 0;
}

.error-text {
  color: #d93025;
  font-size: 12px;
  margin-top: 4px;
}

.button-row {
  dir: ltr;
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-top: 10px;
}

.btn {
  padding: 8px 28px;
  font-size: 14px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  color: white;
}

.btn-save {
  background-color: #28a745;
}

.btn-save:hover:not(:disabled) {
  background-color: #218838;
}

.btn-save:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-clear {
  background-color: #9e9e9e;
}

.btn-clear:hover {
  background-color: #828282;
}

.status-banner {
  text-align: center;
  padding: 10px;
  border-radius: 4px;
  font-size: 14px;
  font-weight: bold;
}

.status-success {
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.status-error {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

@media (max-width: 600px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
  .full-width {
    grid-column: span 1;
  }
  .empty-col {
    display: none;
  }
}
</style>
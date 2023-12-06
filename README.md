# Open Source HRMS (Human Resource Management System) & Employee Portal/Feed

## Project Overview

This project aims to develop an open-source HRMS that not only streamlines HR operations but also focuses on enhancing employee engagement through an interactive and fun employee portal. The core HRMS will manage essential HR tasks while the employee portal will include features such as a radio, an internal social feed, gamification, and more.

### Core HRMS Features Checklist:

- [ ] **Employee Information Management**
  - [ ] Employee profiles
  - [ ] Document management
- [ ] **Leave and Attendance Management**
  - [ ] Leave requests and approval
  - [ ] Attendance tracking
  - [ ] Integration with attendance hardware
- [ ] **Payroll Management**
  - [ ] Salary configuration
  - [ ] Automated payslip generation
  - [ ] Tax calculation
- [ ] **Recruitment Management**
  - [ ] Job posting system
  - [ ] Application tracking
  - [ ] Onboarding workflow
- [ ] **Performance Management**
  - [ ] Appraisal cycles
  - [ ] Goal tracking
  - [ ] Feedback management
- [ ] **Training and Development**
  - [ ] Training scheduling
  - [ ] Skill tracking
- [ ] **Reporting and Analytics**
  - [ ] Standard HR reports
  - [ ] Custom report builder
- [ ] **Benefits Administration**
  - [ ] Insurance and retirement plan management
  - [ ] Employee self-service
- [ ] **Compliance Management**
  - [ ] Compliance document repository
  - [ ] Legal and compliance alerts

### Employee Portal Features Checklist:

- [ ] **Social Feed**
  - [ ] Central newsfeed
  - [ ] Interaction capabilities (post, like, comment)
  - [ ] Events, Polls, etc.
  - [ ] Anonymous posting - to feed or to management
- [ ] **Music Player**
  - [ ] In-built player
  - [ ] Playlist creation and sharing
- [ ] **Gamification and Rewards**
  - [ ] Points and badges
  - [ ] Leaderboards
- [ ] **Health and Wellness Section**
  - [ ] Health situation report and planning
  - [ ] Mental Health and Stress Check-Ins
- [ ] **Personalized Content and Learning**
  - [ ] Internal library access
- [ ] **Recognition Platform**
  - [ ] Peer recognition
- [ ] **Feedback and Survey Tool**
  - [ ] Polls and surveys
- [ ] **Customization and Personalization**
  - [ ] Interface customization
- [ ] **Integration Capabilities**
  - [ ] API integrations
- [ ] **Mobile Application**
  - [ ] Mobile app OR PWA

### Security and Accessibility Features:

- [ ] **Role-Based Access Control (RBAC)**
- [ ] **Single Sign-On (SSO)**
- [ ] **Data Encryption**
- [ ] **Compliance Standards**
- [ ] **Responsive Design**

### ToDo

- [ ] **Access/Refresh JWT Tokens (currently it's 72hours expiry)**

## Contributing

We welcome contributions from the community! If you'd like to contribute, please:

1. Fork the repository.
2. Create a branch for your feature (`git checkout -b feature/amazing-feature`).
3. Commit your changes (`git commit -am 'Add some amazing feature'`).
4. Push to the branch (`git push origin feature/amazing-feature`).
5. Open a pull request.

Please make sure to update the checklist with a `[x]` when you start working on a feature to indicate that it's in progress.

## Development Setup

#### Prerequisites

- Ensure Docker and Docker Compose are installed. [Docker Installation Guide](#)

#### Setup Instructions

1. Navigate to the root of the project folder.
2. To build and start the development environment, run:

   ```bash
   docker-compose build
   docker-compose up

   ```

3. To shut down the development environment, run:

   ```bash
   docker-compose down --remove-orphans

   ```

## Usage

After the containers have been set up, the application FE can be accessed at [http://localhost:3000]. Since the project is in development, the env file is included in the repo and secure features have been commented out.

## Support

## License

This project is licensed under GNU General Public License (GPL) - see the LICENSE file for details.

---

Thank you for your interest in our Open Source HRMS project!

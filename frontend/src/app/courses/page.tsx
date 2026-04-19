'use client';

import { useState, useEffect } from 'react';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { Button } from '@/components/ui/Button';
import { Course } from '@/lib/auth';
import api from '@/lib/api';
import { 
  BookOpen, 
  Clock, 
  DollarSign, 
  Play,
  CheckCircle,
  Lock
} from 'lucide-react';

export default function CoursesPage() {
  const [courses, setCourses] = useState<Course[]>([]);
  const [userCourses, setUserCourses] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'catalog' | 'my-courses'>('catalog');

  useEffect(() => {
    fetchCourses();
    if (activeTab === 'my-courses') {
      fetchUserCourses();
    }
  }, [activeTab]);

  const fetchCourses = async () => {
    try {
      const response = await api.get('/courses/');
      setCourses(response.data.courses || []);
    } catch (error) {
      console.error('Failed to fetch courses:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchUserCourses = async () => {
    try {
      const response = await api.get('/courses/my');
      setUserCourses(response.data.courses || []);
    } catch (error) {
      console.error('Failed to fetch user courses:', error);
    }
  };

  const handleEnroll = async (courseId: string) => {
    try {
      await api.post(`/courses/${courseId}/enroll`);
      fetchUserCourses();
      setActiveTab('my-courses');
    } catch (error: any) {
      console.error('Failed to enroll:', error);
    }
  };

  if (loading) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
        </div>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div>
          <h1 className="text-3xl font-bold text-white mb-2">Courses</h1>
          <p className="text-gray-400">Enhance your trading skills with our comprehensive courses</p>
        </div>

        {/* Tabs */}
        <div className="flex space-x-4 border-b border-gray-700">
          <button
            onClick={() => setActiveTab('catalog')}
            className={`pb-2 px-1 border-b-2 transition-colors ${
              activeTab === 'catalog'
                ? 'border-blue-500 text-blue-400'
                : 'border-transparent text-gray-400 hover:text-gray-300'
            }`}
          >
            Course Catalog
          </button>
          <button
            onClick={() => setActiveTab('my-courses')}
            className={`pb-2 px-1 border-b-2 transition-colors ${
              activeTab === 'my-courses'
                ? 'border-blue-500 text-blue-400'
                : 'border-transparent text-gray-400 hover:text-gray-300'
            }`}
          >
            My Courses
          </button>
        </div>

        {/* Course Catalog */}
        {activeTab === 'catalog' && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {courses.map((course) => (
              <div key={course.id} className="bg-gray-800 rounded-xl border border-gray-700 overflow-hidden">
                {/* Course Image */}
                <div className="h-48 bg-gradient-to-br from-blue-600 to-purple-600 flex items-center justify-center">
                  <BookOpen className="w-16 h-16 text-white/50" />
                </div>

                {/* Course Content */}
                <div className="p-6">
                  <h3 className="text-xl font-semibold text-white mb-2">{course.title}</h3>
                  <p className="text-gray-400 text-sm mb-4 line-clamp-2">{course.description}</p>
                  
                  <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center text-gray-400">
                      <DollarSign className="w-4 h-4 mr-1" />
                      <span className="text-sm">${course.price}</span>
                    </div>
                    <div className="flex items-center text-gray-400">
                      <Clock className="w-4 h-4 mr-1" />
                      <span className="text-sm">Self-paced</span>
                    </div>
                  </div>

                  <Button 
                    className="w-full"
                    onClick={() => handleEnroll(course.id)}
                  >
                    <Play className="w-4 h-4 mr-2" />
                    Enroll Now
                  </Button>
                </div>
              </div>
            ))}

            {courses.length === 0 && (
              <div className="col-span-full text-center py-12">
                <p className="text-gray-400">No courses available at the moment</p>
              </div>
            )}
          </div>
        )}

        {/* My Courses */}
        {activeTab === 'my-courses' && (
          <div className="space-y-6">
            {userCourses.map((userCourse) => (
              <div key={userCourse.id} className="bg-gray-800 p-6 rounded-xl border border-gray-700">
                <div className="flex justify-between items-start mb-4">
                  <div>
                    <h3 className="text-xl font-semibold text-white mb-2">
                      {userCourse.course?.title || 'Course Title'}
                    </h3>
                    <p className="text-gray-400 mb-4">
                      {userCourse.course?.description || 'Course description'}
                    </p>
                    
                    <div className="flex items-center space-x-6 text-sm text-gray-400">
                      <div className="flex items-center">
                        <CheckCircle className="w-4 h-4 mr-1 text-green-400" />
                        <span>Progress: {userCourse.progress?.toFixed(1) || 0}%</span>
                      </div>
                      <div className="flex items-center">
                        {userCourse.completed ? (
                          <CheckCircle className="w-4 h-4 mr-1 text-green-400" />
                        ) : (
                          <Lock className="w-4 h-4 mr-1 text-yellow-400" />
                        )}
                        <span>{userCourse.completed ? 'Completed' : 'In Progress'}</span>
                      </div>
                    </div>
                  </div>

                  <Button>
                    <Play className="w-4 h-4 mr-2" />
                    Continue
                  </Button>
                </div>

                {/* Progress Bar */}
                <div className="w-full bg-gray-700 rounded-full h-2">
                  <div 
                    className="bg-blue-500 h-2 rounded-full transition-all duration-300"
                    style={{ width: `${userCourse.progress || 0}%` }}
                  ></div>
                </div>
              </div>
            ))}

            {userCourses.length === 0 && (
              <div className="text-center py-12">
                <p className="text-gray-400 mb-4">You haven't enrolled in any courses yet</p>
                <Button onClick={() => setActiveTab('catalog')}>
                  Browse Course Catalog
                </Button>
              </div>
            )}
          </div>
        )}
      </div>
    </DashboardLayout>
  );
}

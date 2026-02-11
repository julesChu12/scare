/**
 * WGS84 坐标系（GPS）转 GCJ-02 坐标系（高德/腾讯）
 *
 * 浏览器 Geolocation API 返回的是 WGS84 坐标
 * 高德地图使用的是 GCJ-02 坐标系
 */

const PI = Math.PI
const A = 6378245.0 // 长半轴
const EE = 0.00669342162296594323 // 扁率

/**
 * 判断是否在中国境内
 */
function outOfChina(lng: number, lat: number): boolean {
  return lng < 72.004 || lng > 137.8347 || lat < 0.8293 || lat > 55.8271
}

function transformLat(lng: number, lat: number): number {
  let ret = -100.0 + 2.0 * lng + 3.0 * lat + 0.2 * lat * lat + 0.1 * lng * lat + 0.2 * Math.sqrt(Math.abs(lng))
  ret += (20.0 * Math.sin(6.0 * lng * PI) + 20.0 * Math.sin(2.0 * lng * PI)) * 2.0 / 3.0
  ret += (20.0 * Math.sin(lat * PI) + 40.0 * Math.sin(lat / 3.0 * PI)) * 2.0 / 3.0
  ret += (160.0 * Math.sin(lat / 12.0 * PI) + 320 * Math.sin(lat * PI / 30.0)) * 2.0 / 3.0
  return ret
}

function transformLng(lng: number, lat: number): number {
  let ret = 300.0 + lng + 2.0 * lat + 0.1 * lng * lng + 0.1 * lng * lat + 0.1 * Math.sqrt(Math.abs(lng))
  ret += (20.0 * Math.sin(6.0 * lng * PI) + 20.0 * Math.sin(2.0 * lng * PI)) * 2.0 / 3.0
  ret += (20.0 * Math.sin(lng * PI) + 40.0 * Math.sin(lng / 3.0 * PI)) * 2.0 / 3.0
  ret += (150.0 * Math.sin(lng / 12.0 * PI) + 300.0 * Math.sin(lng / 30.0 * PI)) * 2.0 / 3.0
  return ret
}

/**
 * WGS84 转 GCJ-02
 * @param lng WGS84 经度
 * @param lat WGS84 纬度
 * @returns GCJ-02 坐标 { lng, lat }
 */
export function wgs84ToGcj02(lng: number, lat: number): { lng: number; lat: number } {
  if (outOfChina(lng, lat)) {
    return { lng, lat }
  }

  let dLat = transformLat(lng - 105.0, lat - 35.0)
  let dLng = transformLng(lng - 105.0, lat - 35.0)
  const radLat = lat / 180.0 * PI
  let magic = Math.sin(radLat)
  magic = 1 - EE * magic * magic
  const sqrtMagic = Math.sqrt(magic)
  dLat = (dLat * 180.0) / ((A * (1 - EE)) / (magic * sqrtMagic) * PI)
  dLng = (dLng * 180.0) / (A / sqrtMagic * Math.cos(radLat) * PI)

  return {
    lng: lng + dLng,
    lat: lat + dLat
  }
}

/**
 * 获取当前位置（已转换为 GCJ-02）
 * @returns Promise<{ lng: number; lat: number }> GCJ-02 坐标
 */
export function getCurrentPosition(): Promise<{ lng: number; lat: number }> {
  return new Promise((resolve, reject) => {
    if (!navigator.geolocation) {
      reject(new Error('浏览器不支持地理位置'))
      return
    }

    navigator.geolocation.getCurrentPosition(
      (position) => {
        const { longitude, latitude } = position.coords
        // 转换为 GCJ-02 坐标系
        const gcj02 = wgs84ToGcj02(longitude, latitude)
        resolve(gcj02)
      },
      (error) => {
        let message = '获取位置失败'
        switch (error.code) {
          case error.PERMISSION_DENIED:
            message = '用户拒绝了位置请求'
            break
          case error.POSITION_UNAVAILABLE:
            message = '位置信息不可用'
            break
          case error.TIMEOUT:
            message = '获取位置超时'
            break
        }
        reject(new Error(message))
      },
      {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 60000
      }
    )
  })
}

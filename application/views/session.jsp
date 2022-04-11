<%@page import="java.util.Enumeration"%>
<%@ page language="java" contentType="text/html; charset=EUC-KR"
    pageEncoding="EUC-KR"%>
<!DOCTYPE html>
<html>
<head>
<meta charset="EUC-KR">
<title>Session Get</title>
</head>
<body>
<h1>sessionGet.jsp 입니다.</h1>
<%
    Object sessionObj = session.getAttribute("id");
    String id = sessionObj.toString();
    out.println(id + "<br>");
    
    Object sessionObj2 = session.getAttribute("pw");
    Integer sessionPw = (Integer)sessionObj2;
    out.println(sessionPw + "<br>");
    
    out.println("<hr>");
    
    String sessionName;
    String sessionValue;
    Enumeration enumeration = session.getAttributeNames(); // 세션객체를 직렬화해서 받음
    while(enumeration.hasMoreElements()){
        sessionName = enumeration.nextElement().toString();
        sessionValue = session.getAttribute(sessionName).toString();
        out.println("sessionName : " + sessionName + "<br>");
        out.println("sessionValue : " + sessionValue + "<br>");
    }
    
    out.println("<hr>");
    
    String sessionID = session.getId();    // 브라우저당 session 고유 ID
    out.println("sessionID : " + sessionID + "<br>");
    int s_interval = session.getMaxInactiveInterval();    // session 최대 유효시간 => 톰캣 xml파일에 기본값 60분
    out.println("sessionInterval : " + s_interval + "<br>");
    
    out.println("<hr>");
    
    // 해당 session Name 삭제
    session.removeAttribute("id");
    // 삭제 후 다시 직렬화해서 재 확인
    enumeration = session.getAttributeNames();
    while(enumeration.hasMoreElements()){
        sessionName = enumeration.nextElement().toString();
        sessionValue = session.getAttribute(sessionName).toString();
        out.println("sessionName : " + sessionName + "<br>");
        out.println("sessionValue : " + sessionValue + "<br>");
    }
    
    out.println("<hr>");
    
    session.invalidate();    // session 모두 삭제
    // session이 유효한지 확인을 위해 재요청
    if(request.isRequestedSessionIdValid()) { 
        out.println("session valid");
    }else{
        out.println("session invalid");
    }
%>
</body>
</html>


출처: https://chrismare.tistory.com/entry/세션session [미래를 설계하는 개발자]
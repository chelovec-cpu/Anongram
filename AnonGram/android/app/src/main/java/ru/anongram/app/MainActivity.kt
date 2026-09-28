package ru.anongram.app
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

data class Chat(val title:String,val preview:String,val time:String,val unread:Int)
class MainActivity:ComponentActivity(){override fun onCreate(b:Bundle?){super.onCreate(b);setContent{AnonGramApp()}}}
@Composable fun AnonGramApp(){val chats=listOf(Chat("Анонимный чат","Добро пожаловать в AnonGram","12:40",3),Chat("Hailendsky","Увидимся позже","11:22",0),Chat("CRMP MOBILE","Новый релиз уже близко","09:17",7),Chat("Избранное","Сохранённые сообщения","вчера",0));MaterialTheme(colorScheme=darkColorScheme()){Column(Modifier.fillMaxSize().background(Brush.verticalGradient(listOf(Color(0xFF03040A),Color(0xFF12052A),Color(0xFF03040A))).padding(16.dp))){Text("ANONGRAM",fontSize=27.sp,color=Color(0xFF9C7CFF));Text("Познай анонимность_",fontSize=19.sp,modifier=Modifier.padding(vertical=10.dp));LazyColumn(Modifier.weight(1f)){items(chats){c->ListItem(headlineContent={Text(c.title)},supportingContent={Text(c.preview)},trailingContent={if(c.unread>0)Text(c.unread.toString(),color=Color.Cyan)})}};Button(onClick={},modifier=Modifier.fillMaxWidth()){Text("＋ Новый чат")}}}}
